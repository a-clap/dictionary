package main

import (
	"fmt"
	"github.com/a-clap/dictionary/pkg/translator"
)

func main() {
	s := translator.Translate("brain", translator.English)
	fmt.Println(s)
}
