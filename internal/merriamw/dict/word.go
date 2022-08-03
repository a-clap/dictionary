//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package dict

import (
	"fmt"
	"strings"
	"unicode"
)

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

type Pronunciation struct {
	pron     string
	audioUrl string
}

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
