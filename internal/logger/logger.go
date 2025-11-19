package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() func() error {
	zap := newZapLogger()
	logger = zap.Sugar()
	return zap.Sync
}

func newZapLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"logs/secretaria.log", "stdout"}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
