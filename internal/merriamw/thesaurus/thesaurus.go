//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package thesaurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/logger"
	"io"
	"net/http"
	"net/url"
)

type Thesaurus struct {
	Thesauruser
}

type Thesauruser interface {
	Get(text string) ([]byte, error)
}

type DefaultThesauruser struct {
	key string
}

func NewThesaurus(getWord Thesauruser) *Thesaurus {
	return &Thesaurus{Thesauruser: getWord}
}

func NewThesaurusDefault(key string) *Thesaurus {
	return NewThesaurus(NewDefaultThesauruser(key))
}

func (t *Thesaurus) Translate(text string) (words []*Word, err error) {
	resp, err := t.Get(text)
	if err != nil {
		return nil, fmt.Errorf("error on get %v", err)
	}

	err = json.Unmarshal(resp, &words)
	if err != nil {
		logger.Log.Debugf("error decoding json: %v", err)
		logger.Log.Debugf("parsing as string, to get useful information...")

		var errorInfo []string
		errString := json.Unmarshal(resp, &errorInfo)
		if errString == nil {
			err = fmt.Errorf("%v, additional info %v", err, errorInfo)
		}
		return
	}
	return
}

// Definition simplified access to definition of certain Word
func (w *Word) Definition() []string {
	return w.Shortdef
}

// Text returns word, which translation belongs to
func (w *Word) Text() string {
	return w.Meta.Id
}

// Synonyms simplified access to synonyms
func (w *Word) Synonyms() [][]string {
	return w.Meta.Syns
}

// Antonyms simplified access to antonyms
func (w *Word) Antonyms() [][]string {
	return w.Meta.Ants
}

// IsOffensive returns true, whether word is considered as offensive
func (w *Word) IsOffensive() bool {
	return w.Meta.Offensive
}

// Function returns the word functions in a sentence, e.x. noun, adj etc
func (w *Word) Function() string {
	return w.Fl
}

// NewDefaultThesauruser constructor for default API access
func NewDefaultThesauruser(key string) *DefaultThesauruser {
	return &DefaultThesauruser{key: key}
}

// query returns prepared URL for Get
func (d DefaultThesauruser) query(text string) string {
	const GetUrl = `https://www.dictionaryapi.com/api/v3/references/thesaurus/json/%s?key=%s`
	text = url.PathEscape(text)

	return fmt.Sprintf(GetUrl, text, d.key)

}

// Get fulfills Thesauruser interface
func (d DefaultThesauruser) Get(text string) ([]byte, error) {
	response, err := http.Get(d.query(text))
	if err != nil {
		return nil, fmt.Errorf("get failed %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Debugf("error on Body.Close() %#v", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %s", response.Status)
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(response.Body); err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}

	return buf.Bytes(), nil
}

// Word - structured json, which is received from MerriamWebster, see https://www.dictionaryapi.com/products/json#sec-3
// Deliberately, there are a lot of tags commented - they are not needed by package, however left them as maybe required someday
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
