package logger

import (
	"context"
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// variables
var logger *zap.Logger
var env string

// ContextKey defines keys to be used to pass information to context
type ContextKey string

const (
	DEBUG   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
	FATAL   = 4

	CtxValueRequestId ContextKey = "request-id"
)

// Init initials the logger in various modes as per the environment.
func Init(mode int) {
	var logLevel zapcore.Level
	switch mode {
	case DEBUG:
		logLevel = zapcore.DebugLevel
	case INFO:
		logLevel = zapcore.InfoLevel
	case WARNING:
		logLevel = zapcore.WarnLevel
	case ERROR:
		logLevel = zapcore.ErrorLevel
	case FATAL:
		logLevel = zapcore.FatalLevel
	}

	cfg := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(logLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, _ = cfg.Build()

	if env == "development" {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = logger.WithOptions(
			zap.Fields(zap.Int("pid", os.Getpid()),
				zap.String("exe", path.Base(os.Args[0]))),
			zap.WrapCore(
				func(zapcore.Core) zapcore.Core {
					return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig),
						zapcore.AddSync(os.Stderr), zapcore.DebugLevel)
				}))
	} else {
		logger = logger.WithOptions(
			zap.Fields(zap.Int("pid", os.Getpid()),
				zap.String("exe", path.Base(os.Args[0]))),
			zap.WrapCore(
				func(zapcore.Core) zapcore.Core {
					return zapcore.NewCore(zapcore.NewJSONEncoder(cfg.EncoderConfig), zapcore.AddSync(os.Stderr), logLevel)
				}))
	}
}

// addFieldsFromContext is used to add fields data from context.
func addFieldsFromContext(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	if ctx != nil {
		keyMapping := map[ContextKey]string{
			CtxValueRequestId: "reqId",
		}
		for contextKey, loggerFieldKey := range keyMapping {
			if value, ok := ctx.Value(contextKey).(string); ok {
				fields = append(fields, Field(loggerFieldKey, value))
			}
		}
	}
	return fields
}

func Field(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

func I(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Info(message, addFieldsFromContext(ctx, fields...)...)
}

func D(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Debug(message, addFieldsFromContext(ctx, fields...)...)
}

func W(ctx context.Context, message string, fields ...zapcore.Field) {
	logger.Warn(message, addFieldsFromContext(ctx, fields...)...)
}

func E(ctx context.Context, err error, message string, fields ...zapcore.Field) {
	fields = append(fields, Field("error", err))
	fields = addFieldsFromContext(ctx, fields...)
	logger.Error(message, fields...)
}

func Sync() {
	logger.Info("syncing logger...")
	err := logger.Sync()

	if err != nil {
		fmt.Println("failed syncing logger...")
	}
}
