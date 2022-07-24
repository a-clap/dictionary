package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int64

const (
	Info Level = iota
	Warn
	Debug
	Error
	Fatal
	Panic
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	SetLevel(level Level)
}

type Dummy struct {
}

func (d Dummy) Errorf(string, ...interface{}) {}
func (d Dummy) Fatalf(string, ...interface{}) {}
func (d Dummy) Infof(string, ...interface{})  {}
func (d Dummy) Warnf(string, ...interface{})  {}
func (d Dummy) Debugf(string, ...interface{}) {}
func (d Dummy) SetLevel(Level)                {}

type Standard struct {
	*zap.SugaredLogger
}

func NewDummy() Dummy {
	return Dummy{}
}

func NewDevelopment() Standard {
	//atom := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	logger, _ := cfg.Build()
	return Standard{
		SugaredLogger: logger.Sugar(),
	}
}

func matchLevel(level Level) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	case Panic:
		return zapcore.PanicLevel
	}

	return zapcore.PanicLevel
}

func (s *Standard) SetLevel(level Level) {
	lvl := matchLevel(level)
	s.SugaredLogger = s.SugaredLogger.Desugar().
		WithOptions(zap.IncreaseLevel(lvl)).
		Sugar()
}
