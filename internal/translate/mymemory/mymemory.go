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
	Text  string `json:"translatedText"`
	Match int    `json:"match"`
}

type response struct {
	Data           responseData `json:"responseData"`
	ResponseStatus int          `json:"responseStatus"`
}

func New(logger logger.Logger) *MyMemory {
	return &MyMemory{logger}
}

func getQuery(word string, lang translate.Language) string {
	const getUrl = "https://api.mymemory.translated.net/get?q="
	word = url.QueryEscape(word)
	langPair := "langpair="

	if lang == translate.English {
		langPair += "en|pl"
	} else {
		langPair += "pl|en"
	}
	return fmt.Sprintf("%s%s&%s", getUrl, word, langPair)

}

func (d *MyMemory) Translate(word string, lang translate.Language) (translation string) {
	s := getQuery(word, lang)
	d.Infof("query = %s, for word = %s, lang = %+v", s, word, lang)

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
