package dict

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type GetWord interface {
	Get(text string) ([]byte, error)
}

type Default struct {
	key string
}

func NewDefault(key string) *Default {
	return &Default{key: key}
}

func (d *Default) query(text string) string {
	const GetUrl = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

	text = url.PathEscape(text)
	return fmt.Sprintf(GetUrl, text, d.key)
}

func (d *Default) Get(text string) ([]byte, error) {
	response, err := http.Get(d.query(text))
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
