package utils

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func InitLogger(env string) error {
	var err error
	once.Do(func() {
		if env == "production" {
			logger, err = zap.NewProduction()
		} else {
			logger, err = zap.NewDevelopment()
		}
	})
	return err
}

func SyncLogger() {
	if logger != nil {
		logger.Sync()
	}
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func Errorf(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

type Field = zap.Field

func ErrorField(key string) Field {
	return zap.String("error", key)
}

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func GetLogger() *zap.Logger {
	return logger
}
