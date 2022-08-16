package main

import (
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/pkg/server"
)

func main() {
	s := server.New(nil, logger.NewDevelopment())
	panic(s.Run(":8080"))
}
