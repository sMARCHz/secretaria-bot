package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

func NewProductionLogger() Logger {
	logger := NewZapProduction()
	return NewWithZap(logger)
}

func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

func NewZapProduction() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.NewProductionConfig()
	config.EncoderConfig = encoderConfig
	config.OutputPaths = []string{"logs/secretaria.log", "stderr"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
