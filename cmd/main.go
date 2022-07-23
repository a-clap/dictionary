package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/translate/deepl"
	"os"
)

func main() {
	apikey, found := os.LookupEnv("DEEPL_API_KEY")
	if !found {
		panic("DEEPL_API_KEY not found in ENV")
	}

	d := deepl.New(apikey)
	fmt.Println(d)
}
