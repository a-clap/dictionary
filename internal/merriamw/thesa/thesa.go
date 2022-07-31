package thesa

import (
	"encoding/json"
	"fmt"
)

type Thesaurus struct {
	GetWord
}

func NewThesaurus(getWord GetWord) *Thesaurus {
	return &Thesaurus{GetWord: getWord}
}

func NewThesaurusDefault(key string) *Thesaurus {
	return &Thesaurus{NewDefault(key)}
}

func (t *Thesaurus) Translate(text string) (words []*Word, err error) {
	resp, err := t.Get(text)
	if err != nil {
		return nil, fmt.Errorf("error on get %v", err)
	}

	err = json.Unmarshal(resp, &words)
	if err != nil {
		log.Debugf("error decoding json: %v", err)
		log.Debugf("parsing as string, to get useful information...")

		var errorInfo []string
		errString := json.Unmarshal(resp, &errorInfo)
		if errString == nil {
			err = fmt.Errorf("%v, additional info %v", err, errorInfo)
		}
		return
	}
	return
}
