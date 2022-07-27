package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/merriamw"
	"github.com/a-clap/dictionary/internal/translate"
	"github.com/a-clap/dictionary/internal/translate/mymemory"
	"log"
	"os"
)

func main() {
	thKey, ok := os.LookupEnv("MW_TH_KEY")
	if !ok {
		log.Fatalln("MW_TH_KEY not defined")
	}

	l := logger.NewDevelopment()
	thesaurus := merriamw.NewThesaurus(thKey, &l)

	t, err := thesaurus.Translate("good")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("def =", t.ShortDef())
	fmt.Println("synonyms =", t.Synonyms())
	fmt.Println("antonyms = ", t.Antonyms())

	dictTest()

	mymemtest()
}

func mymemtest() {
	l := logger.NewDevelopment()
	mem := mymemory.New(&l)
	brain := mem.Translate("mózg", translate.Polish)
	fmt.Println(brain)

}

func dictTest() {
	dictKey, ok := os.LookupEnv("MW_DICT_KEY")
	if !ok {
		log.Fatalln("MW_DICT_KEY not defined in env")
	}

	l := logger.NewDevelopment()
	dict := merriamw.NewDictionary(dictKey, &l)

	brain := "brain"
	data, err := dict.Translate(brain)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("shortdef for", brain, "are:\r\n", data.ShortDef())
	fmt.Println(data.AudioFiles())
}
