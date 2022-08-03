//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
)

type DeepL struct {
	Access
	logger.Logger
}

func NewDeepL(access Access, logger logger.Logger) *DeepL {
	return &DeepL{Access: access, Logger: logger}
}

func NewDeepLDefault(key string, logger logger.Logger) *DeepL {
	return NewDeepL(NewAccessDefault(key, logger), logger)
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
