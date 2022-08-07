//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
)

type SourceLang string
type TargetLang string

const (
	SrcBulgarian  SourceLang = "BG"
	SrcCzech      SourceLang = "CS"
	SrcDanish     SourceLang = "DA"
	SrcGerman     SourceLang = "DE"
	SrcGreek      SourceLang = "EL"
	SrcEnglish    SourceLang = "EN"
	SrcSpanish    SourceLang = "ES"
	SrcEstonian   SourceLang = "ET"
	SrcFinnish    SourceLang = "FI"
	SrcFrench     SourceLang = "FR"
	SrcHungarian  SourceLang = "HU"
	SrcIndonesian SourceLang = "ID"
	SrcItalian    SourceLang = "IT"
	SrcJapanese   SourceLang = "JA"
	SrcLithuanian SourceLang = "LT"
	SrcLatvian    SourceLang = "LV"
	SrcDutch      SourceLang = "NL"
	SrcPolish     SourceLang = "PL"
	SrcPortuguese SourceLang = "PT"
	SrcRomanian   SourceLang = "RO"
	SrcRussian    SourceLang = "RU"
	SrcSlovak     SourceLang = "SK"
	SrcSlovenian  SourceLang = "SL"
	SrcSwedish    SourceLang = "SV"
	SrcTurkish    SourceLang = "TR"
	SrcChinese    SourceLang = "ZH"
)
const (
	TarBulgarian       TargetLang = "BG"
	TarCzech           TargetLang = "CS"
	TarDanish          TargetLang = "DA"
	TarGerman          TargetLang = "DE"
	TarGreek           TargetLang = "EL"
	TarEnglishBritish  TargetLang = "EN-GB"
	TarEnglishAmerican TargetLang = "EN-US"
	TarSpanish         TargetLang = "ES"
	TarEstonian        TargetLang = "ET"
	TarFinnish         TargetLang = "FI"
	TarFrench          TargetLang = "FR"
	TarHungarian       TargetLang = "HU"
	TarIndonesian      TargetLang = "ID"
	TarItalian         TargetLang = "IT"
	TarJapanese        TargetLang = "JA"
	TarLithuanian      TargetLang = "LT"
	TarLatvian         TargetLang = "LV"
	TarDutch           TargetLang = "NL"
	TarPolish          TargetLang = "PL"
	TarPortuguese      TargetLang = "PT-PT"
	TarBrazilian       TargetLang = "PT-BR"
	TarRomanian        TargetLang = "RO"
	TarRussian         TargetLang = "RU"
	TarSlovak          TargetLang = "SK"
	TarSlovenian       TargetLang = "SL"
	TarSwedish         TargetLang = "SV"
	TarTurkish         TargetLang = "TR"
	TarChinese         TargetLang = "ZH"
)

type DeepL struct {
	Deepler
	logger.Logger
}

type Deepler interface {
	Query(text string, sourceLang SourceLang, targetLanguage TargetLang) ([]byte, error)
}

func NewDeepL(deepler Deepler, logger logger.Logger) *DeepL {
	return &DeepL{Deepler: deepler, Logger: logger}
}

func NewDeepLDefault(key string, logger logger.Logger) *DeepL {
	return NewDeepL(NewDeeplerDefault(key, logger), logger)
}

func (d *DeepL) Translate(text string, sourceLang SourceLang, targetLang TargetLang) (*Word, error) {
	b, err := d.Query(text, sourceLang, targetLang)
	if err != nil {
		return nil, fmt.Errorf("on query %w", err)
	}
	d.Infof("attempting to parse json")

	w := &Word{}
	err = json.Unmarshal(b, w)
	if err != nil {
		d.Errorf("failed to parse json %#v", err)
		d.Infof("string from data %s", string(b))
		return nil, fmt.Errorf("failed to parse json %w", err)
	}
	return w, nil
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

type DeeplerDefault struct {
	values url.Values
	logger.Logger
}

type Translations struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}
type Word struct {
	Translations []Translations `json:"translations"`
}

func NewDeeplerDefault(key string, logger logger.Logger) *DeeplerDefault {
	return &DeeplerDefault{
		values: map[string][]string{
			"auth_key": {key},
		},
		Logger: logger,
	}
}

func (a *DeeplerDefault) Query(text string, sourceLang SourceLang, targetLang TargetLang) ([]byte, error) {
	a.values.Set("text", text)
	a.values.Set("source_lang", string(sourceLang))
	a.values.Set("target_lang", string(targetLang))

	resp, err := http.PostForm("https://api-free.deepl.com/v2/translate", a.values)
	if err != nil {
		return nil, fmt.Errorf("error on http.Post: %w", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	n, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading response body: %w", err)
	}
	a.Infof("read %v bytes from resp.Body", n)

	return buf.Bytes(), nil
}
