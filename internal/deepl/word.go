//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl

type Word struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func (w Word) SourceLang() []string {
	s := make([]string, len(w.Translations))
	for i, elem := range w.Translations {
		s[i] = elem.DetectedSourceLanguage
	}
	return s
}

func (w Word) Text() []string {
	s := make([]string, len(w.Translations))
	for i, elem := range w.Translations {
		s[i] = elem.Text
	}
	return s
}
