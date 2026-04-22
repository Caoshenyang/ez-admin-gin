package main

import (
	"ez-admin-gin/server/internal/bootstrap"
	"ez-admin-gin/server/internal/permission"

	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"
	"ez-admin-gin/server/internal/token"

	"go.uber.org/zap"
)

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

	// 数据库连接成功后，创建基础表并准备默认管理员。
	if err := bootstrap.Run(db, log); err != nil {
		log.Fatal("bootstrap application", zap.Error(err))
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
