package main

import (
	// stdlog 只用于日志系统初始化失败前的兜底输出。
	stdlog "log"

	"ez-admin-gin/server/internal/config"
	appLogger "ez-admin-gin/server/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 先读取配置，日志初始化也需要用到 cfg.Log。
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
		// 退出前刷新缓冲区，避免最后几条日志丢失。
		_ = log.Sync()
	}()

	// 使用 gin.New()，再手动挂载自定义中间件。
	r := gin.New()
	r.Use(appLogger.GinLogger(log), appLogger.GinRecovery(log))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.App.Env,
		})
	})

	// 服务启动日志记录关键运行参数。
	log.Info(
		"server started",
		zap.String("addr", cfg.Server.Addr),
		zap.String("env", cfg.App.Env),
	)

	if err := r.Run(cfg.Server.Addr); err != nil {
		// Fatal 会记录日志并退出进程。
		log.Fatal("run server", zap.Error(err))
	}
}
