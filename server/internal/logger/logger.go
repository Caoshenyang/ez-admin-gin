package logger

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ez-admin-gin/server/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New 根据配置创建 Zap Logger。
// 这里同时处理日志级别、日志格式和输出位置。
func New(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// 生产配置的字段更稳定，适合后续接入日志采集平台。
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	var encoder zapcore.Encoder
	if strings.EqualFold(cfg.Format, "json") {
		// json 适合生产环境采集和检索。
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// console 适合本地开发时直接阅读。
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	writeSyncer, err := newWriteSyncer(cfg)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)

	return zap.New(
		core,
		// 记录调用位置，方便定位日志来自哪个文件。
		zap.AddCaller(),
		// error 级别自动带堆栈，方便排查异常。
		zap.AddStacktrace(zapcore.ErrorLevel),
	), nil
}

// GinLogger 记录每一次 HTTP 请求。
func GinLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 先执行后续处理器，才能拿到最终状态码和耗时。
		c.Next()

		if query != "" {
			path = path + "?" + query
		}

		// 用结构化字段记录请求信息，后续按字段过滤会更方便。
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			// 5xx 请求按错误日志记录。
			log.Error("http request", fields...)
			return
		}

		log.Info("http request", fields...)
	}
}

// GinRecovery 捕获 panic，避免服务因为单次请求异常直接退出。
func GinRecovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error(
			"panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.Stack("stack"),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// parseLevel 把配置文件中的字符串转成 Zap 识别的日志级别。
func parseLevel(value string) (zapcore.Level, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(value))); err != nil {
		return level, fmt.Errorf("invalid log level %q: %w", value, err)
	}

	return level, nil
}

// newWriteSyncer 决定日志输出到哪里。
func newWriteSyncer(cfg config.LogConfig) (zapcore.WriteSyncer, error) {
	if cfg.Filename == "" {
		// 没有配置文件路径时，只输出到控制台。
		return zapcore.AddSync(os.Stdout), nil
	}

	// 日志目录不存在时自动创建，例如 logs/app.log 会先创建 logs/。
	if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0o755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	// Lumberjack 负责日志切割，避免单个日志文件无限增长。
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 同时输出到控制台和文件：开发时能直接看，事后也能查文件。
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(fileWriter),
	), nil
}
