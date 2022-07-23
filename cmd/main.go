package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/translate"
	"github.com/a-clap/dictionary/internal/translate/mymemory"
)

func main() {
	l := logger.NewDevelopment()
	mem := mymemory.New(l)
	s := mem.Translate("hello", translate.English)
	fmt.Println(s)
}
