package dict

import (
	"encoding/json"
	"fmt"
)

type Dict struct {
	GetWord
}

func NewDictDefault(key string) *Dict {
	return NewDict(NewDefault(key))
}

func NewDict(getWord GetWord) *Dict {
	return &Dict{
		GetWord: getWord,
	}
}

func (d *Dict) Translate(text string) (data []*Word, err error) {
	resp, err := d.Get(text)
	if err != nil {
		log.Errorf("error on get %v", err)
		return
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Debugf("error decoding json: %v", err)
		log.Debugf("parsing as string, to get useful information...")

		var errorInfo []string
		errString := json.Unmarshal(resp, &errorInfo)
		if errString == nil {
			err = fmt.Errorf("%v, additional info %v", err, errorInfo)
		}
		return
	}
	return
}
