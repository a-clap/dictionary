package translator

import (
	"github.com/a-clap/dictionary/internal/logger"
	dictionary "github.com/a-clap/dictionary/internal/merriamw/dict"
	thesaurus "github.com/a-clap/dictionary/internal/merriamw/thesa"
	"github.com/a-clap/dictionary/internal/mymemory"
	"os"
)

type Language mymemory.Language

const (
	Polish  Language = Language(mymemory.Polish)
	English Language = Language(mymemory.English)
)

// Translation contains everything, what translator is able to give
type Translation struct {
	Text        string   `json:"text"`
	Translation string   `json:"translation"`
	Antonyms    []string `json:"antonyms"`
	Synonyms    []string `json:"synonyms"`
	IsOffensive bool     `json:"isOffensive"`
	Audio       []string `json:"audio"`
}

// Suggestions if there is typo in text (or translator is not able to translator), in Suggestions will be stored slice of strings with possible texts
type Suggestions struct {
	Suggestions []string
}

var (
	translator struct {
		dict *dictionary.Dict
		thes *thesaurus.Thesaurus
		mw   *mymemory.MyMemory
	}
)

// init get needed API keys from ENV variables, then initializes translator
func init() {
	log := logger.NewDevelopment()
	dictKey, ok := os.LookupEnv("MW_DICT_KEY")
	if !ok {
		log.Errorf("MW_DICT_KEY not defined in ENV")
		os.Exit(-1)
	}

	thKey, ok := os.LookupEnv("MW_TH_KEY")
	if !ok {
		log.Errorf("MW_TH_KEY not defined in ENV")
		os.Exit(-1)
	}

	translator.dict = dictionary.NewDictDefault(dictKey, log)
	translator.thes = thesaurus.NewThesaurusDefault(thKey, log)
	translator.mw = mymemory.NewMyMemoryDefault(log)
}

func Translate(text string, from Language) Translation {

	return Translation{
		Translation: "",
		Antonyms:    nil,
		Synonyms:    nil,
		IsOffensive: false,
		Audio:       nil,
	}
}
