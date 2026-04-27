package main

import (
	"ez-admin-gin/server/internal/permission"

	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appMigrate "ez-admin-gin/server/internal/migrate"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"
	"ez-admin-gin/server/internal/token"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	// 嵌入迁移文件
	"embed"
)

//go:embed migrations/pgsql migrations/mysql
var migrationsFS embed.FS

func main() {
	// 先读取配置，日志、数据库、Redis 初始化都依赖配置。
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	// 根据配置创建结构化日志对象。
	log, err := appLogger.New(cfg.Log)
	if err != nil {
		stdlog.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	// 启动时连接数据库；连接失败就直接终止服务。
	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Error("close database", zap.Error(err))
		}
	}()

	// 数据库连接成功后，执行 SQL 迁移（建表 + 种子数据）。
	migrateDSN, err := database.DSN(cfg.Database)
	if err != nil {
		log.Fatal("build migration dsn", zap.Error(err))
	}
	if cfg.Database.Driver == "postgres" {
		migrateDSN = cfg.Database.Driver + "://" + migrateDSN
	}
	if err := appMigrate.Run(cfg.Database.Driver, migrateDSN, migrationsFS, log); err != nil {
		if err == migrate.ErrNoChange {
			log.Info("database migrations up to date")
		} else {
			log.Fatal("run database migrations", zap.Error(err))
		}
	}

	// 启动时连接 Redis；连接失败就直接终止服务。
	redisClient, err := appRedis.New(cfg.Redis, log)
	if err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}
	defer func() {
		if err := appRedis.Close(redisClient); err != nil {
			log.Error("close redis", zap.Error(err))
		}
	}()

	// Token 管理器负责签发和解析登录令牌。
	tokenManager, err := token.NewManager(cfg.Auth)
	if err != nil {
		log.Fatal("create token manager", zap.Error(err))
	}

	// 权限判断器负责根据角色策略判断接口访问权限。
	permissionEnforcer, err := permission.NewEnforcer(db, "configs/rbac_model.conf")
	if err != nil {
		log.Fatal("create permission enforcer", zap.Error(err))
	}

	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{
		Config:     cfg,
		Log:        log,
		DB:         db,
		Redis:      redisClient,
		Token:      tokenManager,
		Permission: permissionEnforcer,
	})

	// 服务启动日志记录关键运行参数。
	log.Info(
		"server started",
		zap.String("addr", cfg.Server.Addr),
		zap.String("env", cfg.App.Env),
	)

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatal("run server", zap.Error(err))
	}
}
