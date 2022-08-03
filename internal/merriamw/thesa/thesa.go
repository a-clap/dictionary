package thesa

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
)

type Thesaurus struct {
	GetWord
	logger.Logger
}

func NewThesaurus(getWord GetWord, logger logger.Logger) *Thesaurus {
	return &Thesaurus{GetWord: getWord, Logger: logger}
}

func NewThesaurusDefault(key string, logger logger.Logger) *Thesaurus {
	return NewThesaurus(NewDefault(key), logger)
}

func (t *Thesaurus) Translate(text string) (words []*Word, err error) {
	resp, err := t.Get(text)
	if err != nil {
		return nil, fmt.Errorf("error on get %v", err)
	}

	err = json.Unmarshal(resp, &words)
	if err != nil {
		t.Debugf("error decoding json: %v", err)
		t.Debugf("parsing as string, to get useful information...")

		var errorInfo []string
		errString := json.Unmarshal(resp, &errorInfo)
		if errString == nil {
			err = fmt.Errorf("%v, additional info %v", err, errorInfo)
		}
		return
	}
	return
}
