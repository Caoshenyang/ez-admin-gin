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

	"time"

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
	Engine    *gin.Engine
	Config    *platformConfig.Config
	DB        *gorm.DB
	Redis     *goredis.Client
	Token     *authnPlatform.Manager
	Enforcer  *authzPlatform.Enforcer
	Log       *zap.Logger
	Server    *httptest.Server
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
	sessionStore := authnPlatform.NewRedisSessionStore(rdb, 24*time.Hour)

	engine := bootstrap.NewRouter(bootstrap.RouterOptions{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      rdb,
		Token:      token,
		Session:    sessionStore,
		Permission: enforcer,
	})

	srv := httptest.NewServer(engine)

	return &TestApp{
		Engine:   engine,
		Config:   cfg,
		DB:       db,
		Redis:    rdb,
		Token:    token,
		Enforcer: enforcer,
		Log:      log,
		Server:   srv,
	}
}

// CleanupTestData removes test-generated departments, roles, users, and casbin rules
// that are not part of the seed data. This enables test isolation without
// requiring a full database reset between tests.
func (app *TestApp) CleanupTestData(t *testing.T) {
	t.Helper()

	// Delete casbin rules for non-super_admin roles.
	app.DB.Exec("DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 != 'super_admin'")
	// Delete role-data-scope bindings for non-seed roles.
	app.DB.Exec("DELETE FROM sys_role_data_scope WHERE role_id > 1")
	// Delete role-menu bindings for non-seed roles.
	app.DB.Exec("DELETE FROM sys_role_menu WHERE role_id > 1")
	// Delete user-role bindings for non-seed users (ID > 1).
	app.DB.Exec("DELETE FROM sys_user_role WHERE user_id > 1")
	// Delete user-post bindings for non-seed users.
	app.DB.Exec("DELETE FROM sys_user_post WHERE user_id > 1")
	// Delete roles that are not the built-in super_admin (ID = 1).
	app.DB.Exec("DELETE FROM sys_role WHERE id > 1")
	// Delete users that are not the seed admin (ID = 1).
	app.DB.Exec("DELETE FROM sys_user WHERE id > 1")
	// Delete all departments (none in seed data).
	app.DB.Exec("DELETE FROM sys_department")
	// Delete all menus (none in seed data).
	app.DB.Exec("DELETE FROM sys_menu")

	// Reload policies to clear stale entries from memory.
	app.ReloadPolicies(t)
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

// SeedRestrictedUser creates a test role with no permissions and a test user
// assigned to that role. Returns the user's access token.
//
// This enables RBAC tests where a user has limited (or zero) API permissions.
func (app *TestApp) SeedRestrictedUser(t *testing.T, adminToken, username, password, nickname string) string {
	t.Helper()

	// 1. Create a role with no permissions via API.
	roleBody := `{"code":"test_no_perm","name":"Test No Permission","sort":99,"data_scope":"self","status":1}`
	roleReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", adminToken, strings.NewReader(roleBody))
	roleResp, err := app.Do(roleReq)
	if err != nil {
		t.Fatalf("create role request failed: %v", err)
	}
	defer roleResp.Body.Close()

	if roleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(roleResp.Body)
		t.Fatalf("create role failed: status %d, body: %s", roleResp.StatusCode, string(body))
	}

	var roleResult struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatalf("decode role response: %v", err)
	}
	roleID := roleResult.Data.ID
	if roleID == 0 {
		t.Fatal("created role has ID 0")
	}

	// 2. Create a user assigned to the new role via API.
	userBody := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s","department_id":0,"status":1,"role_ids":[%d]}`,
		username, password, nickname, roleID,
	)
	userReq := app.AuthRequest(http.MethodPost, "/api/v1/system/users", adminToken, strings.NewReader(userBody))
	userResp, err := app.Do(userReq)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(userResp.Body)
		t.Fatalf("create user failed: status %d, body: %s", userResp.StatusCode, string(body))
	}

	// 3. Login as the new user and return the token.
	return app.LoginAs(t, username, password)
}

// SeedUserWithPermissions creates a test role with specific API permissions and
// a test user assigned to that role. Returns the user's access token.
//
// permissions is a list of {path, method} pairs that will be granted to the role.
func (app *TestApp) SeedUserWithPermissions(t *testing.T, adminToken, username, password, nickname, roleCode, roleName string, permissions [][2]string) string {
	t.Helper()

	// 1. Create a role.
	roleBody := fmt.Sprintf(
		`{"code":"%s","name":"%s","sort":98,"data_scope":"self","status":1}`,
		roleCode, roleName,
	)
	roleReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", adminToken, strings.NewReader(roleBody))
	roleResp, err := app.Do(roleReq)
	if err != nil {
		t.Fatalf("create role request failed: %v", err)
	}
	defer roleResp.Body.Close()

	if roleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(roleResp.Body)
		t.Fatalf("create role failed: status %d, body: %s", roleResp.StatusCode, string(body))
	}

	var roleResult struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatalf("decode role response: %v", err)
	}
	roleID := roleResult.Data.ID
	if roleID == 0 {
		t.Fatal("created role has ID 0")
	}

	// 2. Assign permissions to the role.
	permItems := make([]string, 0, len(permissions))
	for _, p := range permissions {
		permItems = append(permItems, fmt.Sprintf(`{"path":"%s","method":"%s"}`, p[0], p[1]))
	}
	permBody := fmt.Sprintf(`{"permissions":[%s]}`, strings.Join(permItems, ","))
	permReq := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/permissions", roleID), adminToken, strings.NewReader(permBody))
	permResp, err := app.Do(permReq)
	if err != nil {
		t.Fatalf("update permissions request failed: %v", err)
	}
	defer permResp.Body.Close()

	if permResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(permResp.Body)
		t.Fatalf("update permissions failed: status %d, body: %s", permResp.StatusCode, string(body))
	}

	// 3. Create a user assigned to the role.
	userBody := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s","department_id":0,"status":1,"role_ids":[%d]}`,
		username, password, nickname, roleID,
	)
	userReq := app.AuthRequest(http.MethodPost, "/api/v1/system/users", adminToken, strings.NewReader(userBody))
	userResp, err := app.Do(userReq)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(userResp.Body)
		t.Fatalf("create user failed: status %d, body: %s", userResp.StatusCode, string(body))
	}

	// 4. Login as the new user.
	return app.LoginAs(t, username, password)
}

// SeedDepartment creates a single department via API. Returns the department ID.
func (app *TestApp) SeedDepartment(t *testing.T, adminToken string, parentID uint, name, code string) uint {
	t.Helper()

	body := fmt.Sprintf(
		`{"parent_id":%d,"name":"%s","code":"%s","sort":0,"status":1}`,
		parentID, name, code,
	)
	req := app.AuthRequest(http.MethodPost, "/api/v1/system/departments", adminToken, strings.NewReader(body))
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("create department request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("create department failed: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode department response: %v", err)
	}
	if result.Data.ID == 0 {
		t.Fatal("created department has ID 0")
	}
	return result.Data.ID
}

// SeedScopedUser creates a role with a specific data_scope and a user in a
// given department. The role also gets the specified API permissions.
// Returns the user's access token.
func (app *TestApp) SeedScopedUser(t *testing.T, adminToken, username, password, nickname, roleCode, roleName string, dataScope string, departmentID uint, permissions [][2]string) string {
	t.Helper()

	// 1. Create a role with the specified data_scope.
	roleBody := fmt.Sprintf(
		`{"code":"%s","name":"%s","sort":97,"data_scope":"%s","status":1}`,
		roleCode, roleName, dataScope,
	)
	roleReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", adminToken, strings.NewReader(roleBody))
	roleResp, err := app.Do(roleReq)
	if err != nil {
		t.Fatalf("create role request failed: %v", err)
	}
	defer roleResp.Body.Close()

	if roleResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(roleResp.Body)
		t.Fatalf("create role failed: status %d, body: %s", roleResp.StatusCode, string(respBody))
	}

	var roleResult struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatalf("decode role response: %v", err)
	}
	roleID := roleResult.Data.ID
	if roleID == 0 {
		t.Fatal("created role has ID 0")
	}

	// 2. Assign permissions to the role.
	if len(permissions) > 0 {
		permItems := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permItems = append(permItems, fmt.Sprintf(`{"path":"%s","method":"%s"}`, p[0], p[1]))
		}
		permBody := fmt.Sprintf(`{"permissions":[%s]}`, strings.Join(permItems, ","))
		permReq := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/permissions", roleID), adminToken, strings.NewReader(permBody))
		permResp, err := app.Do(permReq)
		if err != nil {
			t.Fatalf("update permissions request failed: %v", err)
		}
		defer permResp.Body.Close()

		if permResp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(permResp.Body)
			t.Fatalf("update permissions failed: status %d, body: %s", permResp.StatusCode, string(respBody))
		}
	}

	// 3. Create a user in the specified department with this role.
	userBody := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s","department_id":%d,"status":1,"role_ids":[%d]}`,
		username, password, nickname, departmentID, roleID,
	)
	userReq := app.AuthRequest(http.MethodPost, "/api/v1/system/users", adminToken, strings.NewReader(userBody))
	userResp, err := app.Do(userReq)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(userResp.Body)
		t.Fatalf("create user failed: status %d, body: %s", userResp.StatusCode, string(respBody))
	}

	// 4. Login as the new user.
	return app.LoginAs(t, username, password)
}

// SeedCustomDeptUser creates a role with data_scope=custom_dept and explicitly specified
// department IDs, then creates a user assigned to that role. Returns the user's access token.
func (app *TestApp) SeedCustomDeptUser(t *testing.T, adminToken, username, password, nickname, roleCode, roleName string, customDeptIDs []uint, userDeptID uint, permissions [][2]string) string {
	t.Helper()

	// 1. Build custom_department_ids JSON array.
	deptIDStrs := make([]string, len(customDeptIDs))
	for i, id := range customDeptIDs {
		deptIDStrs[i] = fmt.Sprintf("%d", id)
	}
	roleBody := fmt.Sprintf(
		`{"code":"%s","name":"%s","sort":96,"data_scope":"custom_dept","custom_department_ids":[%s],"status":1}`,
		roleCode, roleName, strings.Join(deptIDStrs, ","),
	)
	roleReq := app.AuthRequest(http.MethodPost, "/api/v1/system/roles", adminToken, strings.NewReader(roleBody))
	roleResp, err := app.Do(roleReq)
	if err != nil {
		t.Fatalf("create custom_dept role request failed: %v", err)
	}
	defer roleResp.Body.Close()

	if roleResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(roleResp.Body)
		t.Fatalf("create custom_dept role failed: status %d, body: %s", roleResp.StatusCode, string(respBody))
	}

	var roleResult struct {
		Code int `json:"code"`
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatalf("decode custom_dept role response: %v", err)
	}
	roleID := roleResult.Data.ID
	if roleID == 0 {
		t.Fatal("created custom_dept role has ID 0")
	}

	// 2. Assign permissions to the role.
	if len(permissions) > 0 {
		permItems := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permItems = append(permItems, fmt.Sprintf(`{"path":"%s","method":"%s"}`, p[0], p[1]))
		}
		permBody := fmt.Sprintf(`{"permissions":[%s]}`, strings.Join(permItems, ","))
		permReq := app.AuthRequest(http.MethodPost, fmt.Sprintf("/api/v1/system/roles/%d/permissions", roleID), adminToken, strings.NewReader(permBody))
		permResp, err := app.Do(permReq)
		if err != nil {
			t.Fatalf("update custom_dept role permissions request failed: %v", err)
		}
		defer permResp.Body.Close()

		if permResp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(permResp.Body)
			t.Fatalf("update custom_dept role permissions failed: status %d, body: %s", permResp.StatusCode, string(respBody))
		}
	}

	// 3. Create a user in the specified department with this role.
	userBody := fmt.Sprintf(
		`{"username":"%s","password":"%s","nickname":"%s","department_id":%d,"status":1,"role_ids":[%d]}`,
		username, password, nickname, userDeptID, roleID,
	)
	userReq := app.AuthRequest(http.MethodPost, "/api/v1/system/users", adminToken, strings.NewReader(userBody))
	userResp, err := app.Do(userReq)
	if err != nil {
		t.Fatalf("create custom_dept user request failed: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(userResp.Body)
		t.Fatalf("create custom_dept user failed: status %d, body: %s", userResp.StatusCode, string(respBody))
	}

	// 4. Login as the new user.
	return app.LoginAs(t, username, password)
}

// DecodeResponse decodes a JSON response body into the provided target.
func (app *TestApp) DecodeResponse(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

// ReloadPolicies reloads Casbin policies from the database into the in-memory
// enforcer. Call this after creating roles or updating permissions via the API
// so that subsequent requests use the latest policy state.
func (app *TestApp) ReloadPolicies(t *testing.T) {
	t.Helper()
	if err := app.Enforcer.ReloadPolicy(); err != nil {
		t.Fatalf("reload casbin policies: %v", err)
	}
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
