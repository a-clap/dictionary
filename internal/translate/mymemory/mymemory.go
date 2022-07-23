package mymemory

import "github.com/a-clap/dictionary/internal/translate"

type MyMemory struct {
	apiKey string
}

func New(apiKey string) *MyMemory {
	return &MyMemory{apiKey: apiKey}
}

func (d *MyMemory) Translate(word string, lang translate.Language) (translation string) {
	return ""
}
