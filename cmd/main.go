package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/deepl"
	"github.com/a-clap/dictionary/internal/logger"
	"os"
)

func main() {
	log := logger.NewDevelopment()
	key, ok := os.LookupEnv("DEEPL_KEY")
	if !ok {
		log.Fatalf("DEEPL_KEY not found in env")
	}
	d := deepl.NewDeepLDefault(key, log)
	w, err := d.Translate("i like dumplings", deepl.SrcEnglish, deepl.TarPolish)
	if err != nil {
		log.Fatalf("%#v", err)
	}
	log.Infof("parsed successfully\n")
	fmt.Println(w.Text())
}
