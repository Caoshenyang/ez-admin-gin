// Package testutil provides test infrastructure for centralized integration tests.
//
// Tests that require PostgreSQL or Redis should use the "integration" build tag:
//
//	//go:build integration
//
// and be run via:
//
//	go test -tags integration ./server/tests/api/...
package testutil

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"ez-admin-gin/server/internal/bootstrap"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"
	platformLogger "ez-admin-gin/server/internal/platform/logger"
	platformMigrate "ez-admin-gin/server/internal/platform/migrate"
	platformRedis "ez-admin-gin/server/internal/platform/redis"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:embed testdata/config.yaml
var testdataFS embed.FS

// TestApp holds a fully-bootstrapped test application and its dependencies.
type TestApp struct {
	Engine *gin.Engine
	Config *platformConfig.Config
	DB     *gorm.DB
	Redis  *goredis.Client
	Token  *authnPlatform.Manager
	Log    *zap.Logger
	Server *httptest.Server
}

// NewTestApp bootstraps a complete test application.
//
// It loads config from testdata/config.yaml (env vars with EZ_ prefix still
// take precedence), connects to a real PostgreSQL and Redis, runs migrations,
// and returns an httptest.Server wrapping the Gin engine.
//
// Callers must defer app.Close() to release resources.
func NewTestApp(t *testing.T) *TestApp {
	t.Helper()

	gin.SetMode(gin.TestMode)

	cfg := mustLoadTestConfig(t)
	log := mustCreateLogger(t, cfg)
	db := mustConnectDB(t, cfg, log)
	mustRunMigrations(t, cfg, log)
	rdb := mustConnectRedis(t, cfg, log)
	token := mustCreateTokenManager(t, cfg)
	enforcer := mustCreateEnforcer(t, db)

	engine := bootstrap.NewRouter(bootstrap.RouterOptions{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      rdb,
		Token:      token,
		Permission: enforcer,
	})

	srv := httptest.NewServer(engine)

	return &TestApp{
		Engine: engine,
		Config: cfg,
		DB:     db,
		Redis:  rdb,
		Token:  token,
		Log:    log,
		Server: srv,
	}
}

// Close releases all resources held by the test application.
func (app *TestApp) Close(t *testing.T) {
	t.Helper()
	app.Server.Close()
	if err := platformRedis.Close(app.Redis); err != nil {
		t.Logf("close redis: %v", err)
	}
	if err := platformDatabase.Close(app.DB); err != nil {
		t.Logf("close database: %v", err)
	}
}

// SeedAdmin calls the setup/init endpoint to create an initial admin user.
// This should be called once before tests that need a logged-in user.
func (app *TestApp) SeedAdmin(t *testing.T, username, password, nickname string) {
	t.Helper()

	body := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s"}`,
		username, password, nickname,
	)
	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/setup/init"), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("setup/init request failed: %v", err)
	}
	defer resp.Body.Close()

	// 200 or 409 (already initialized) are both acceptable.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("setup/init failed: status %d, body: %s", resp.StatusCode, string(body))
	}
}

// LoginAs logs in with the given credentials and returns the access token.
func (app *TestApp) LoginAs(t *testing.T, username, password string) string {
	t.Helper()

	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/login"), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: status %d", resp.StatusCode)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if result.Data.AccessToken == "" {
		t.Fatal("login returned empty access_token")
	}
	return result.Data.AccessToken
}

// AuthRequest wraps an HTTP request with a Bearer token header.
func (app *TestApp) AuthRequest(method, path, token string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, app.URL(path), body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Do sends an HTTP request to the test server and returns the response.
func (app *TestApp) Do(req *http.Request) (*http.Response, error) {
	return app.Server.Client().Do(req)
}

// URL returns the full test server URL for a given path.
func (app *TestApp) URL(path string) string {
	return app.Server.URL + path
}

// NewRequest creates a new HTTP request targeting the test server.
func (app *TestApp) NewRequest(method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, app.URL(path), body)
	if err != nil {
		panic(err)
	}
	return req
}

// --- internal helpers ---

func mustLoadTestConfig(t *testing.T) *platformConfig.Config {
	t.Helper()

	v := viper.New()
	v.SetConfigType("yaml")

	data, err := testdataFS.ReadFile("testdata/config.yaml")
	if err != nil {
		t.Fatalf("read test config: %v", err)
	}

	v.ReadConfig(strings.NewReader(string(data)))

	// Allow env vars to override test config.
	v.SetEnvPrefix("EZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindTestEnvs(v)

	var cfg platformConfig.Config
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatalf("unmarshal test config: %v", err)
	}
	return &cfg
}

func mustCreateLogger(t *testing.T, cfg *platformConfig.Config) *zap.Logger {
	t.Helper()
	log, err := platformLogger.New(cfg.Log)
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	return log
}

func mustConnectDB(t *testing.T, cfg *platformConfig.Config, log *zap.Logger) *gorm.DB {
	t.Helper()
	db, err := platformDatabase.New(cfg.Database, log)
	if err != nil {
		t.Skipf("database not available, skipping integration test: %v", err)
	}
	return db
}

func mustRunMigrations(t *testing.T, cfg *platformConfig.Config, log *zap.Logger) {
	t.Helper()
	migrateDSN, err := platformDatabase.MigrateDSN(cfg.Database)
	if err != nil {
		t.Fatalf("build migration dsn: %v", err)
	}

	if err := platformMigrate.Run(cfg.Database.Driver, migrateDSN, migrationsOnly(), log); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
}

func mustConnectRedis(t *testing.T, cfg *platformConfig.Config, log *zap.Logger) *goredis.Client {
	t.Helper()
	rdb, err := platformRedis.New(cfg.Redis, log)
	if err != nil {
		t.Skipf("redis not available, skipping integration test: %v", err)
	}
	return rdb
}

func mustCreateTokenManager(t *testing.T, cfg *platformConfig.Config) *authnPlatform.Manager {
	t.Helper()
	mgr, err := authnPlatform.NewManager(cfg.Auth)
	if err != nil {
		t.Fatalf("create token manager: %v", err)
	}
	return mgr
}

func mustCreateEnforcer(t *testing.T, db *gorm.DB) *authzPlatform.Enforcer {
	t.Helper()
	enforcer, err := authzPlatform.NewEnforcer(db, "../../configs/rbac_model.conf")
	if err != nil {
		t.Skipf("create enforcer (db may not have casbin tables yet): %v", err)
	}
	return enforcer
}

func migrationsOnly() fs.FS {
	// migrate.Run expects the FS to contain "migrations/postgres/" at root level.
	// We point to the server/ directory which has "migrations/" as a subdirectory.
	return os.DirFS("../..")
}

func bindTestEnvs(v *viper.Viper) {
	keys := []string{
		"database.host", "database.port", "database.user", "database.password", "database.name",
		"redis.host", "redis.port", "redis.password", "redis.db",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
