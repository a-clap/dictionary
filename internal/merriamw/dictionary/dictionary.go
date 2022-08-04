//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package dictionary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

type Dictionary struct {
	GetWord
	logger.Logger
}

type GetWord interface {
	Get(text string) ([]byte, error)
}

type Pronunciation struct {
	PhoneticNotation string
	Url              string
}

type DefaultGetWord struct {
	key string
}

type Suggestions struct {
	Suggestions []string
}

func NewDictDefault(key string, logger logger.Logger) *Dictionary {
	return NewDictionary(NewDefaultGetWord(key), logger)
}

func NewDictionary(getWord GetWord, logger logger.Logger) *Dictionary {
	return &Dictionary{
		GetWord: getWord,
		Logger:  logger,
	}
}

// Definition return possible slice of Definition to text.
func (d Dictionary) Definition(text string) (data []*Definition, err error) {
	resp, err := d.Get(text)
	if err != nil {
		err = fmt.Errorf("error on get %w", err)
		d.Errorf("error on get %v", err)
		return
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		// This usually means, text wasn't found on dictionary.
		// In that case, we will get an array of strings with suggestions
		d.Debugf("error decoding json: %v", err)
		d.Debugf("parsing as string, to get useful information...")

		var errorInfo []string
		errString := json.Unmarshal(resp, &errorInfo)
		if errString == nil {
			err = fmt.Errorf("%v, additional info %v", err, errorInfo)
		}
		return
	}
	return
}

// Definition simplified access to definition of certain word
func (w *Definition) Definition() []string {
	return w.Shortdef
}

// Text returns word, which translation belongs to
func (w *Definition) Text() string {
	// MerriamW sometimes adds unique number to each word after ":". We can get this unique number and associate words as homographs
	return strings.Split(w.Meta.Id, ":")[0]
}

// Examples returns slice of strings with usage of certain word
func (w *Definition) Examples() []string {
	examples := make([]string, len(w.Suppl.Examples))
	for i, elem := range w.Suppl.Examples {
		examples[i] = elem.T
	}
	return examples
}

// Audio returns possible pronunciations for word
func (w *Definition) Audio() []Pronunciation {
	const AudioUrl = `https://media.merriam-webster.com/audio/prons/en/us/mp3/%s/%s.mp3`

	prons := make([]Pronunciation, 0, len(w.Hwi.Prs))
	for _, elem := range w.Hwi.Prs {
		pron := Pronunciation{
			PhoneticNotation: elem.Mw,
		}

		filename := elem.Sound.Audio
		if len(filename) > 0 {
			dir := ""
			if strings.HasPrefix(filename, "bix") {
				dir = "bix"
			} else if strings.HasPrefix(filename, "gg") {
				dir = "gg"
			} else if unicode.IsNumber(rune(filename[0])) || unicode.IsPunct(rune(filename[0])) {
				dir = "number"
			} else {
				dir = string(filename[0])
			}
			pron.Url = fmt.Sprintf(AudioUrl, dir, filename)
		}
		prons = append(prons, pron)
	}
	return prons
}

// IsOffensive returns true, whether word is considered as offensive
func (w *Definition) IsOffensive() bool {
	return w.Meta.Offensive
}

// Function returns the word functions in a sentence, e.x. noun, adj etc
func (w *Definition) Function() string {
	return w.Fl
}

// NewDefaultGetWord constructor for standard API access
func NewDefaultGetWord(key string) *DefaultGetWord {
	return &DefaultGetWord{key: key}
}

// query returns prepared URL for Get
func (d DefaultGetWord) query(text string) string {
	const GetUrl = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

	text = url.PathEscape(text)
	return fmt.Sprintf(GetUrl, text, d.key)
}

// Get fulfills GetWord interface
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

// Definition - structured json, which is received from MerriamWebster, see https://www.dictionaryapi.com/products/json#sec-2
// Deliberately, there are a lot of tags commented - they are not needed by package, however left them as maybe required someday
type Definition struct {
	Meta struct {
		Id   string `json:"id"`
		Uuid string `json:"uuid"`
		//Sort      string   `json:"sort"`
		//Src       string   `json:"src"`
		//Section   string   `json:"section"`
		//Stems     []string `json:"stems"`
		Offensive bool `json:"offensive"`
	} `json:"meta"`
	Hwi struct {
		Hw  string `json:"hw"`
		Prs []struct {
			Mw    string `json:"mw"`
			Sound struct {
				Audio string `json:"audio"`
				Ref   string `json:"ref"`
				Stat  string `json:"stat"`
			} `json:"sound"`
		} `json:"prs"`
	} `json:"hwi"`
	Fl string `json:"fl"`
	//Def []struct {
	//	Sseq [][][]interface{} `json:"sseq"`
	//} `json:"def"`
	//Uros []struct {
	//	Ure string `json:"ure"`
	//	Fl  string `json:"fl"`
	//} `json:"uros"`
	//Et     [][]string `json:"et"`
	//Date   string     `json:"date"`
	//LdLink struct {
	//	LinkHw string `json:"link_hw"`
	//	LinkFl string `json:"link_fl"`
	//} `json:"ld_link"`
	Suppl struct {
		Examples []struct {
			T string `json:"t"`
		} `json:"examples"`
		//	Ldq struct {
		//		Ldhw string `json:"ldhw"`
		//		Fl   string `json:"fl"`
		//		Def  []struct {
		//			Sls  []string          `json:"sls"`
		//			Sseq [][][]interface{} `json:"sseq"`
		//		} `json:"def"`
		//	} `json:"ldq"`
	} `json:"suppl"`
	Shortdef []string `json:"shortdef"`
}
