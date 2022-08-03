//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int64

// Level defines current logging level
const (
	Info Level = iota
	Warn
	Debug
	Error
	Fatal
	Panic
)

// Logger interface
type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	SetLevel(level Level)
}

// Dummy satisfies Logger interface, it just doesn't log anything anywhere
type Dummy struct {
}

func (d Dummy) Errorf(string, ...interface{}) {}
func (d Dummy) Fatalf(string, ...interface{}) {}
func (d Dummy) Infof(string, ...interface{})  {}
func (d Dummy) Warnf(string, ...interface{})  {}
func (d Dummy) Debugf(string, ...interface{}) {}
func (d Dummy) SetLevel(Level)                {}

// NewDummy creates Dummy Logger
func NewDummy() Dummy {
	return Dummy{}
}

// Development logger type, inherits from zap.SugaredLogger
type Development struct {
	*zap.SugaredLogger
}

// NewDevelopment create new Development logger with predefined time layout
func NewDevelopment() *Development {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	logger, _ := cfg.Build()
	return &Development{
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
	// Just to discard warning
	return zapcore.PanicLevel
}

// SetLevel dynamic set level of logging
func (s *Development) SetLevel(level Level) {
	lvl := matchLevel(level)
	s.SugaredLogger = s.SugaredLogger.Desugar().
		WithOptions(zap.IncreaseLevel(lvl)).
		Sugar()
}
