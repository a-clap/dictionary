package mymemory

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/translate"
)

type MyMemory struct {
	GetWord
}

func NewMyMemory(word GetWord) *MyMemory {
	return &MyMemory{GetWord: word}
}

func NewMyMemoryDefault() *MyMemory {
	return &MyMemory{NewDefault()}
}

func (d *MyMemory) Translate(word string, lang translate.Language) (words *Word, err error) {
	data, err := d.Get(word, lang)
	if err != nil {
		return nil, fmt.Errorf("error on get: %v", err)
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &words)
	if err != nil {
		log.Errorf("error decoding json %v", err)
		return
	}
	return
}
