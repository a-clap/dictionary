//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package mymemory

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/logger"
)

type MyMemory struct {
	GetWord
}

func NewMyMemory(word GetWord) *MyMemory {
	return &MyMemory{GetWord: word}
}

func NewMyMemoryDefault() *MyMemory {
	return NewMyMemory(NewDefault())
}

func (d *MyMemory) Translate(word string, lang Language) (words *Word, err error) {
	data, err := d.Get(word, lang)
	if err != nil {
		return nil, fmt.Errorf("error on get: %v", err)
	}

	err = json.Unmarshal(data, &words)
	if err != nil {
		logger.Errorf("error decoding json %v", err)
		return
	}
	return
}
