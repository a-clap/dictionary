package main

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/merriamw/dict"
	"github.com/a-clap/dictionary/internal/merriamw/thesa"
	"github.com/a-clap/dictionary/internal/translate"
	"github.com/a-clap/dictionary/internal/translate/mymemory"
	"log"
	"os"
)

func main() {
	l := logger.NewDevelopment()
	thesa.InitLogger(l)
	mymemory.InitLogger(l)
	mymemory.InitLogger(l)

	//dictTest()
	//thtest()

	mymemtest()
}

func thtest() {
	thKey, ok := os.LookupEnv("MW_TH_KEY")
	if !ok {
		log.Fatalln("MW_TH_KEY not defined")
	}

	thesaurus := thesa.NewThesaurusDefault(thKey)
	t, err := thesaurus.Translate("face")
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range t {
		fmt.Println("Text:", elem.Text())
		fmt.Println("Shortdef:", elem.Definition())
		fmt.Println("Synonyms:", elem.Synonyms())
		fmt.Println("Antonyms:", elem.Antonyms())
	}

}

func mymemtest() {
	mem := mymemory.NewMyMemoryDefault()

	data, err := mem.Translate("rack one's brain", translate.English)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(data.Translated())
	fmt.Println(data.Alternatives())

}

func dictTest() {
	dictKey, ok := os.LookupEnv("MW_DICT_KEY")
	if !ok {
		log.Fatalln("MW_DICT_KEY not defined in env")
	}

	dictDefault := dict.NewDictDefault(dictKey)

	brain := "get"
	data, err := dictDefault.Translate(brain)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %v translations\n", len(data))
	for _, elem := range data {
		fmt.Printf("Text : %v\n", elem.Text())
		fmt.Printf("Shortdef: %v\n", elem.Definition())
		fmt.Printf("Pronunciation: %v\n", elem.Audio())
	}
}
