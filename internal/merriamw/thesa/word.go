//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package thesa

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
