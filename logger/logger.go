package logger

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
	Sugar *zap.SugaredLogger
)

func Init(env string) error {
	var err error
	if env == "dev" {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}
	if err != nil{
		return err
	}
	Sugar = Logger.Sugar()
	return nil
}
func Sync() {
    _ = Logger.Sync()
}