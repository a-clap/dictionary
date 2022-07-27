package merriamw

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
	"strings"
)

type Thesaurus struct {
	key string
	logger.Logger
}

func NewThesaurus(key string, logger logger.Logger) *Thesaurus {
	return &Thesaurus{key: key, Logger: logger}
}

func (t *Thesaurus) query(text string) string {
	const GetUrl = `https://www.dictionaryapi.com/api/v3/references/thesaurus/json/%s?key=%s`
	text = url.PathEscape(text)

	return fmt.Sprintf(GetUrl, text, t.key)

}

func (t *Thesaurus) Translate(text string) (*ThesaurusWord, error) {
	query := t.query(text)
	t.Infof("query = %s, for text %s", query, text)

	resp, err := http.Get(query)
	if err != nil {
		t.Errorf("error on http get %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("received status different than OK %s", resp.Status)
		return nil, fmt.Errorf("wrong status code %s", resp.Status)
	}

	word := &ThesaurusWord{Words: nil}
	err = json.NewDecoder(resp.Body).Decode(&word.Words)
	if err != nil {
		t.Debugf("error decoding json to ThesaurusWord %v", err)
		return nil, err
	}
	return word, nil
}

func (t *ThesaurusWord) ShortDef() string {
	return strings.Join(t.Words[0].Shortdef, "\n")
}
func (t *ThesaurusWord) Synonyms() string {
	if len(t.Words[0].Meta.Syns) > 0 {
		return strings.Join(t.Words[0].Meta.Syns[0], "\n")
	}
	return ""
}
func (t *ThesaurusWord) Antonyms() string {
	if len(t.Words[0].Meta.Ants) > 0 {
		return strings.Join(t.Words[0].Meta.Ants[0], "\n")
	}
	return ""

}

type ThesaurusWord struct {
	Words []ThesaurusSingleWord
}

type ThesaurusSingleWord struct {
	Meta struct {
		Id      string `json:"id"`
		Uuid    string `json:"uuid"`
		Src     string `json:"src"`
		Section string `json:"section"`
		Target  struct {
			Tuuid string `json:"tuuid"`
			Tsrc  string `json:"tsrc"`
		} `json:"target"`
		Stems     []string   `json:"stems"`
		Syns      [][]string `json:"syns"`
		Ants      [][]string `json:"ants"`
		Offensive bool       `json:"offensive"`
	} `json:"meta"`
	Hwi struct {
		Hw string `json:"hw"`
	} `json:"hwi"`
	Fl  string `json:"fl"`
	Def []struct {
		Sseq [][][]interface{} `json:"sseq"`
	} `json:"def"`
	Shortdef []string `json:"shortdef"`
}
