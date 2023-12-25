package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitProductionLogger() func() error {
	zap := NewZapProduction()
	logger = zap.Sugar()
	return zap.Sync
}

func NewZapProduction() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.NewProductionConfig()
	config.EncoderConfig = encoderConfig
	config.OutputPaths = []string{"logs/secretaria.log", "stdout"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
