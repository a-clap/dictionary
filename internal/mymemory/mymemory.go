//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package mymemory

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
)

type MyMemory struct {
	GetWord
	logger.Logger
}

func NewMyMemory(word GetWord, logger logger.Logger) *MyMemory {
	return &MyMemory{GetWord: word, Logger: logger}
}

func NewMyMemoryDefault(logger logger.Logger) *MyMemory {
	return NewMyMemory(NewDefault(), logger)
}

func (d *MyMemory) Translate(word string, lang Language) (words *Word, err error) {
	data, err := d.Get(word, lang)
	if err != nil {
		return nil, fmt.Errorf("error on get: %v", err)
	}

	err = json.Unmarshal(data, &words)
	if err != nil {
		d.Errorf("error decoding json %v", err)
		return
	}
	return
}
