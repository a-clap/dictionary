package deepl

import "github.com/a-clap/dictionary/internal/translate"

type Deepl struct {
	apiKey string
}

func New(apiKey string) *Deepl {
	return &Deepl{apiKey: apiKey}
}

func (d *Deepl) Translate(word string, lang translate.Language) (translation string) {
	return ""
}
