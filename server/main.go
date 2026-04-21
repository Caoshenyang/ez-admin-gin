package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"
	"net/http"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 先读取配置，日志和数据库初始化都依赖配置。
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

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New()
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))

	r.GET("/health", func(c *gin.Context) {
		if err := database.Ping(db); err != nil {
			log.Error("database health check failed", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"env":      cfg.App.Env,
				"database": "unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"env":      cfg.App.Env,
			"database": "ok",
		})
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
