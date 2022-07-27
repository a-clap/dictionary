package mymemory

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/translate"
	"net/http"
	"net/url"
)

type MyMemory struct {
	logger.Logger
}

type responseData struct {
	Text  string  `json:"translatedText"`
	Match float64 `json:"match"`
}

type response struct {
	Data           responseData `json:"responseData"`
	ResponseStatus int          `json:"responseStatus"`
}

func New(logger logger.Logger) *MyMemory {
	return &MyMemory{logger}
}

func query(word string, lang translate.Language) string {
	const GetUrl = "https://api.mymemory.translated.net/get?q=%s&langpair=%s"
	word = url.PathEscape(word)

	langPair := ""
	if lang == translate.English {
		langPair = "en|pl"
	} else {
		langPair = "pl|en"
	}
	return fmt.Sprintf(GetUrl, word, langPair)

}

func (d *MyMemory) Translate(word string, lang translate.Language) (translation string) {
	s := query(word, lang)
	d.Infof("query = %s, for text = %s, lang = %+v", s, word, lang)

	resp, err := http.Get(s)
	if err != nil {
		d.Errorf("error %v", err)
		return
	}

	var r response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		d.Errorf("error decoding json %v", err)
		return
	}
	d.Infof("decoded json %+v", r)

	return r.Data.Text
}
