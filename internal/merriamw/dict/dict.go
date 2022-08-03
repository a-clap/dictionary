//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package dict

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
)

type Dict struct {
	GetWord
	logger.Logger
}

func NewDictDefault(key string, logger logger.Logger) *Dict {
	return NewDict(NewDefault(key), logger)
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
