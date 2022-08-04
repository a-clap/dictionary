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

type Dict struct {
	GetWord
	logger.Logger
}

type GetWord interface {
	Get(text string) ([]byte, error)
}

type Pronunciation struct {
	pron     string
	audioUrl string
}

type DefaultGetWord struct {
	key string
}

func NewDictDefault(key string, logger logger.Logger) *Dict {
	return NewDict(NewDefaultGetWord(key), logger)
}

func NewDict(getWord GetWord, logger logger.Logger) *Dict {
	return &Dict{
		GetWord: getWord,
		Logger:  logger,
	}
}

func (d *Dict) Translate(text string) (data []*Word, err error) {
	resp, err := d.Get(text)
	if err != nil {
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

func (w *Word) Definition() []string {
	return w.Shortdef
}

func (w *Word) Text() string {
	// MerriamW sometimes adds unique number to each word after ":". We can get this unique number and associate words as homographs
	return strings.Split(w.Meta.Id, ":")[0]
}

func (w *Word) Examples() []string {
	examples := make([]string, 0, len(w.Suppl.Examples))
	for _, elem := range w.Suppl.Examples {
		examples = append(examples, elem.T)
	}
	return examples
}

func (w *Word) Audio() []Pronunciation {
	const AudioUrl = `https://media.merriam-webster.com/audio/prons/en/us/mp3/%s/%s.mp3`

	prons := make([]Pronunciation, 0, len(w.Hwi.Prs))
	for _, elem := range w.Hwi.Prs {
		pron := Pronunciation{
			pron: elem.Mw,
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
			pron.audioUrl = fmt.Sprintf(AudioUrl, dir, filename)
		}
		prons = append(prons, pron)
	}
	return prons
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
	const GetUrl = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

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

// Word - structured json, which is received from MerriamWebster, see https://www.dictionaryapi.com/products/json#sec-2
// Deliberately, there are a lot of tags commented - they are not needed by package, however left them as maybe required someday
type Word struct {
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
