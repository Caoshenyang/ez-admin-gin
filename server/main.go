package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"
	appRedis "ez-admin-gin/server/internal/redis"
	"ez-admin-gin/server/internal/router"

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

	// 路由注册交给 internal/router，main.go 只保留启动流程。
	r := router.New(router.Options{
		Config: cfg,
		Log:    log,
		DB:     db,
		Redis:  redisClient,
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
