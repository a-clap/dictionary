package main

import (
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/a-clap/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger.Init(logger.NewDefaultZap(zapcore.DebugLevel))

	s := server.New(nil)
	panic(s.Run(":8080"))
}
