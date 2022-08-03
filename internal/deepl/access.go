//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl

import (
	"bytes"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
)

type SourceLang string
type TargetLang string

const (
	Bulgarian  SourceLang = "BG"
	Czech                 = "CS"
	Danish                = "DA"
	German                = "DE"
	Greek                 = "EL"
	English               = "EN"
	Spanish               = "ES"
	Estonian              = "ET"
	Finnish               = "FI"
	French                = "FR"
	Hungarian             = "HU"
	Indonesian            = "ID"
	Italian               = "IT"
	Japanese              = "JA"
	Lithuanian            = "LT"
	Latvian               = "LV"
	Dutch                 = "NL"
	Polish                = "PL"
	Portuguese            = "PT"
	Romanian              = "RO"
	Russian               = "RU"
	Slovak                = "SK"
	Slovenian             = "SL"
	Swedish               = "SV"
	Turkish               = "TR"
	Chinese               = "ZH"
)

type Access interface {
	Query(text string, sourceLang SourceLang, targetLanguage TargetLang) ([]byte, error)
}

type AccessDefault struct {
	values url.Values
	logger.Logger
}

func NewAccessDefault(key string, logger logger.Logger) *AccessDefault {
	return &AccessDefault{
		values: map[string][]string{
			"auth_key": {key},
		},
		Logger: logger,
	}
}

func (a *AccessDefault) Query(text string, sourceLang SourceLang, targetLang TargetLang) ([]byte, error) {
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
