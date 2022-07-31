package mymemory

import (
	"bytes"
	"fmt"
	"github.com/a-clap/dictionary/internal/translate"
	"net/http"
	"net/url"
)

type GetWord interface {
	Get(text string, lang translate.Language) ([]byte, error)
}

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func query(text string, lang translate.Language) string {
	const GetUrl = "https://api.mymemory.translated.net/get?q=%s&langpair=%s"
	text = url.PathEscape(text)

	langPair := ""
	if lang == translate.English {
		langPair = "en|pl"
	} else {
		langPair = "pl|en"
	}
	return fmt.Sprintf(GetUrl, text, langPair)
}

func (d *Default) Get(text string, lang translate.Language) ([]byte, error) {
	response, err := http.Get(query(text, lang))
	if err != nil {
		return nil, fmt.Errorf("get failed %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %s", response.Status)
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(response.Body); err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}

	return buf.Bytes(), nil
}
