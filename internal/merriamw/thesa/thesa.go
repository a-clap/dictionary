//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package thesa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
)

type Thesaurus struct {
	GetWord
	logger.Logger
}

type GetWord interface {
	Get(text string) ([]byte, error)
}

type DefaultGetWord struct {
	key string
}

func NewThesaurus(getWord GetWord, logger logger.Logger) *Thesaurus {
	return &Thesaurus{GetWord: getWord, Logger: logger}
}

func NewThesaurusDefault(key string, logger logger.Logger) *Thesaurus {
	return NewThesaurus(NewDefaultGetWord(key), logger)
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

func (w *Word) Definition() []string {
	return w.Shortdef
}

func (w *Word) Text() string {
	return w.Meta.Id
}

func (w *Word) Synonyms() [][]string {
	return w.Meta.Syns
}
func (w *Word) Antonyms() [][]string {
	return w.Meta.Ants
}

func (w *Word) IsOffensive() bool {
	return w.Meta.Offensive
}

func (w *Word) Function() string {
	return w.Fl
}

func NewDefaultGetWord(key string) *DefaultGetWord {
	return &DefaultGetWord{key: key}
}

func (d DefaultGetWord) query(text string) string {
	const GetUrl = `https://www.dictionaryapi.com/api/v3/references/thesaurus/json/%s?key=%s`
	text = url.PathEscape(text)

	return fmt.Sprintf(GetUrl, text, d.key)

}

func (d DefaultGetWord) Get(text string) ([]byte, error) {
	response, err := http.Get(d.query(text))
	if err != nil {
		return nil, fmt.Errorf("get failed %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %s", response.Status)
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(response.Body); err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}

	return buf.Bytes(), nil
}

type Word struct {
	Meta struct {
		Id string `json:"id"`
		//Uuid    string `json:"uuid"`
		//Src     string `json:"src"`
		//Section string `json:"section"`
		//Target  struct {
		//	Tuuid string `json:"tuuid"`
		//	Tsrc  string `json:"tsrc"`
		//} `json:"target"`
		//Stems     []string   `json:"stems"`
		Syns      [][]string `json:"syns"`
		Ants      [][]string `json:"ants"`
		Offensive bool       `json:"offensive"`
	} `json:"meta"`
	//Hwi struct {
	//	Hw string `json:"hw"`
	//} `json:"hwi"`
	Fl string `json:"fl"`
	//Def []struct {
	//	Sseq [][][]interface{} `json:"sseq"`
	//} `json:"def"`
	Shortdef []string `json:"shortdef"`
}
