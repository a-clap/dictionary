//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package dictionary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

var Logger logger.Logger = logger.NewNop()

type Dictionary struct {
	Definitioner
}

type Definitioner interface {
	Get(text string) ([]byte, error)
}

type Pronunciation struct {
	PhoneticNotation string
	Url              string
}

type DefaultGetDefinition struct {
	key string
}

type Suggestions struct {
	Suggestions []string
}

func NewDictDefault(key string) *Dictionary {
	return NewDictionary(NewDefaultGetDefinition(key))
}

func NewDictionary(getDefinition Definitioner) *Dictionary {
	return &Dictionary{
		Definitioner: getDefinition,
	}
}

// Definition return possible slice of Definition for passed argument.
// If it couldn't find exact Definition, function may returned slice with Suggestions - if there is a typo in word.
// Otherwise error
func (d Dictionary) Definition(text string) (data []*Definition, suggestions *Suggestions, err error) {
	resp, err := d.Get(text)
	if err != nil {
		err = fmt.Errorf("error on get %w", err)
		Logger.Errorf("error on get %v", err)
		return
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		data = nil
		// This usually means, text wasn't found on dictionary.
		// In that case, we will get an array of strings with suggestions
		Logger.Debugf("error decoding json: %v", err)
		Logger.Debugf("parsing as string, to get useful information...")

		suggestions = &Suggestions{Suggestions: []string{}}
		errString := json.Unmarshal(resp, &suggestions.Suggestions)
		if errString == nil {
			err = nil
			Logger.Debugf("...success!")
		} else {
			suggestions = nil
			err = fmt.Errorf("%w %v", err, errString)
			Logger.Debugf("...failure!")
		}
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

	prons := make([]Pronunciation, len(w.Hwi.Prs))
	for i, elem := range w.Hwi.Prs {
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
		prons[i] = pron
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

// NewDefaultGetDefinition constructor for standard API access
func NewDefaultGetDefinition(key string) *DefaultGetDefinition {
	return &DefaultGetDefinition{key: key}
}

// query returns prepared URL for Get
func (d DefaultGetDefinition) query(text string) string {
	const GetUrl = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

	text = url.PathEscape(text)
	return fmt.Sprintf(GetUrl, text, d.key)
}

// Get fulfills Definitioner interface
func (d DefaultGetDefinition) Get(text string) ([]byte, error) {
	response, err := http.Get(d.query(text))
	if err != nil {
		return nil, fmt.Errorf("get failed %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Debugf("error on Body.Close() %#v", err)
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
