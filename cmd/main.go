package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/translate/mymemory"
	"os"
)

func main() {
	apikey, found := os.LookupEnv("MY_MEMORY_API_KEY")
	if !found {
		panic("MY_MEMORY_API_KEY not found in ENV")
	}

	d := mymemory.New(apikey)
	fmt.Println(d)
}
