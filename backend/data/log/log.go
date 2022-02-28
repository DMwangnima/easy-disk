package log

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	BaseLogger Logger
)

func init() {
	var err error
	BaseLogger, err = newZapLogger()
	if err != nil {
		panic(fmt.Sprintf("log: initialize logger %s", err.Error()))
	}
}

type Logger interface {
	Debug(template string, args ...interface{})
	Info(template string, args ...interface{})
    Error(template string, args ...interface{})
	Panic(template string, args ...interface{})
	Fatal(template string, args ...interface{})
	With(args ...interface{}) Logger
}

type zapLogger struct {
	internal *zap.SugaredLogger
}

func newZapLogger() (*zapLogger, error) {
	base, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := base.Sugar()
	return &zapLogger{internal: sugar}, nil
}

func (logger *zapLogger) Debug(template string, args ...interface{}) {
	logger.internal.Debugf(template, args...)
}

func (logger *zapLogger) Info(template string, args ...interface{}) {
	logger.internal.Infof(template, args...)
}

func (logger *zapLogger) Error(template string, args ...interface{}) {
	logger.internal.Errorf(template, args...)
}

func (logger *zapLogger) Panic(template string, args ...interface{}) {
	logger.internal.Panicf(template, args...)
}

func (logger *zapLogger) Fatal(template string, args ...interface{}) {
	logger.internal.Fatalf(template, args...)
}

func (logger *zapLogger) With(args ...interface{}) Logger {
	sugar := logger.internal.With(args...)
	return &zapLogger{internal: sugar}
}



