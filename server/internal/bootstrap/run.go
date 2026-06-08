package bootstrap

import (
	"io/fs"
	stdlog "log"
	"net"
	"strings"

	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	authzPlatform "ez-admin-gin/server/internal/platform/authz"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	platformDatabase "ez-admin-gin/server/internal/platform/database"
	platformLogger "ez-admin-gin/server/internal/platform/logger"
	platformMigrate "ez-admin-gin/server/internal/platform/migrate"
	platformRedis "ez-admin-gin/server/internal/platform/redis"

	"go.uber.org/zap"
)

// MustRun 启动后台服务；初始化失败时直接终止进程。
func MustRun(migrationsFS fs.FS, rbacModelPath string) {
	cfg, err := platformConfig.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	if err := cfg.ValidateProduction(); err != nil {
		stdlog.Fatalf("config validation: %v", err)
	}

	log, err := platformLogger.New(cfg.Log)
	if err != nil {
		stdlog.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	db, err := platformDatabase.New(cfg.Database, log)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}
	defer func() {
		if err := platformDatabase.Close(db); err != nil {
			log.Error("close database", zap.Error(err))
		}
	}()

	migrateDSN, err := platformDatabase.MigrateDSN(cfg.Database)
	if err != nil {
		log.Fatal("build migration dsn", zap.Error(err))
	}
	if err := platformMigrate.Run(cfg.Database.Driver, migrateDSN, migrationsFS, log); err != nil {
		log.Fatal("run database migrations", zap.Error(err))
	}

	redisClient, err := platformRedis.New(cfg.Redis, log)
	if err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}
	defer func() {
		if err := platformRedis.Close(redisClient); err != nil {
			log.Error("close redis", zap.Error(err))
		}
	}()

	tokenManager, err := authnPlatform.NewManager(cfg.Auth)
	if err != nil {
		log.Fatal("create token manager", zap.Error(err))
	}

	permissionEnforcer, err := authzPlatform.NewEnforcer(db, rbacModelPath)
	if err != nil {
		log.Fatal("create permission enforcer", zap.Error(err))
	}

	r := NewRouter(RouterOptions{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      redisClient,
		Token:      tokenManager,
		Permission: permissionEnforcer,
	})

	log.Info(
		"server started",
		zap.String("addr", cfg.Server.Addr),
		zap.String("env", cfg.App.Env),
	)
	if cfg.Swagger.Enabled {
		log.Info("swagger ui available", zap.String("url", swaggerURL(cfg.Server.Addr)))
	}

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatal("run server", zap.Error(err))
	}
}

func swaggerURL(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		trimmed := strings.TrimSpace(addr)
		if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
			return strings.TrimRight(trimmed, "/") + "/swagger/index.html"
		}
		return "http://" + strings.TrimRight(trimmed, "/") + "/swagger/index.html"
	}

	if host == "" || host == "0.0.0.0" || host == "::" || host == "[::]" {
		host = "localhost"
	}
	return "http://" + net.JoinHostPort(host, port) + "/swagger/index.html"
}
