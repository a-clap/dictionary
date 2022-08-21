//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package translator

import (
	"github.com/a-clap/dictionary/internal/deepl"
	"github.com/a-clap/dictionary/internal/merriamw/dictionary"
	"github.com/a-clap/dictionary/internal/merriamw/thesaurus"
	"github.com/a-clap/logger"
)

type Translate interface {
	Get(text string, from deepl.SourceLang, to deepl.TargetLang) (*Translation, error)
}

type DeeplTranslate struct {
	Text string `json:"text"`
}
type Definition struct {
	Offensive  bool                       `json:"offensive"`
	Function   string                     `json:"function"`
	Examples   []string                   `json:"examples"`
	Definition []string                   `json:"definition"`
	Audio      []dictionary.Pronunciation `json:"audio"`
}

type DictionaryTranslate struct {
	Defs     []Definition `json:"defs"`
	Synonyms []string     `json:"synonyms"`
}

type ThesaurusTranslate struct {
	Text       string   `json:"text"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
	Offensive  bool     `json:"offensive"`
	Function   string   `json:"function"`
	Definition []string `json:"definition"`
}

// Translation contains everything, what can be received from Translator
type Translation struct {
	Deepl      []DeeplTranslate     `json:"deepl"`
	Dictionary *DictionaryTranslate `json:"dictionary"`
	Thesaurus  []ThesaurusTranslate `json:"thesaurus"`
}

type Translator struct {
	Translate
}

type standard struct {
	deepl     *deepl.DeepL
	dict      *dictionary.Dictionary
	thesaurus *thesaurus.Thesaurus
}

func (s *standard) Get(text string, from deepl.SourceLang, to deepl.TargetLang) (*Translation, error) {
	deeplTranslate, err := s.deepl.Translate(text, from, to)
	if err != nil {
		return nil, err
	}

	t := &Translation{
		Deepl:      make([]DeeplTranslate, len(deeplTranslate.Translations)),
		Dictionary: nil,
		Thesaurus:  nil,
	}

	for i, elem := range deeplTranslate.Translations {
		logger.Log.Infof("got translation %s", elem.Translation())
		t.Deepl[i].Text = elem.Translation()
	}

	t.Dictionary = s.getDefinitions(to, &t.Deepl)
	t.Thesaurus = s.getThesaurus(to, &t.Deepl)
	return t, nil
}

func (s *standard) getDefinitions(to deepl.TargetLang, deeplTranslate *[]DeeplTranslate) *DictionaryTranslate {
	// Currently supported only for english
	if to != deepl.TarEnglishAmerican && to != deepl.TarEnglishBritish {
		return nil
	}

	dictTranslates := &DictionaryTranslate{
		Defs:     []Definition{},
		Synonyms: []string{},
	}

	for _, elem := range *deeplTranslate {
		text := elem.Text
		d, _, err := s.dict.Definition(text)
		if err != nil || d == nil {
			logger.Log.Debugf("definition not found")
			continue
		}

		for _, dict := range d {
			logger.Log.Debugf("definition for %s is %s", text, dict.Text())
			if dict.Text() != text {
				logger.Log.Debugf("skipping definition as it is not equal text, adding as synonym")
				dictTranslates.Synonyms = append(dictTranslates.Synonyms, dict.Text())
				continue
			}
			dictTranslate := Definition{
				Offensive:  dict.IsOffensive(),
				Function:   dict.Function(),
				Examples:   dict.Examples(),
				Definition: dict.Definition(),
				Audio:      dict.Audio(),
			}

			dictTranslates.Defs = append(dictTranslates.Defs, dictTranslate)
		}
	}
	return dictTranslates
}

func (s *standard) getThesaurus(to deepl.TargetLang, deeplTranslates *[]DeeplTranslate) []ThesaurusTranslate {
	// Currently supported only for english
	if to != deepl.TarEnglishAmerican && to != deepl.TarEnglishBritish {
		return nil
	}
	var th []ThesaurusTranslate
	for _, elem := range *deeplTranslates {
		text := elem.Text
		data, err := s.thesaurus.Translate(text)
		if err != nil {
			logger.Log.Debugf("thesaurus not found for text %s", text)
			continue
		}

		for _, elem := range data {
			if elem.Text() != text {
				continue
			}

			t := ThesaurusTranslate{
				Text:       elem.Text(),
				Synonyms:   nil,
				Antonyms:   nil,
				Offensive:  elem.IsOffensive(),
				Function:   elem.Function(),
				Definition: elem.Definition(),
			}

			if len(elem.Synonyms()) > 0 {
				t.Synonyms = elem.Synonyms()[0]
			}

			if len(elem.Antonyms()) > 0 {
				t.Antonyms = elem.Antonyms()[0]
			}

			th = append(th, t)
			// Naive implementation - get just one Thesaurus for each elem.Text()
			break
		}

	}

	return th
}

func New(translate Translate) *Translator {
	return &Translator{Translate: translate}
}

func NewStandard(deeplKey, dictKey, thKey string) *Translator {
	standard := &standard{
		deepl:     deepl.NewDeepLDefault(deeplKey),
		dict:      dictionary.NewDictDefault(dictKey),
		thesaurus: thesaurus.NewThesaurusDefault(thKey),
	}

	return New(standard)
}
