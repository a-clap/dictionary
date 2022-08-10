package main

import "github.com/a-clap/dictionary/pkg/server"

func main() {
	s := server.New(nil)
	s.Run(":8080")
}
