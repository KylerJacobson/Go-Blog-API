package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Sugar() *zap.SugaredLogger
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

func NewLogger(env string) (*zap.Logger, error) {
	if env == "dev" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
