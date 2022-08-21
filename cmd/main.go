package main

import (
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/a-clap/logger"
)

func main() {
	s := server.New(nil, logger.NewDevelopment())
	panic(s.Run(":8080"))
}
