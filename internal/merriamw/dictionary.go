package merriamw

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

type Dictionary struct {
	key string
	logger.Logger
}

type DictionaryWord struct {
	Words  []DictionarySingleWord
	logger logger.Logger
}

func (d *DictionaryWord) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &d.Words)
	if err != nil {
		var s []string
		if err = json.Unmarshal(data, &s); err == nil {
			return fmt.Errorf("%v:\n%s", err, strings.Join(s, "\n"))
		}
	}
	return err
}

func NewDictionary(key string, logger logger.Logger) *Dictionary {
	return &Dictionary{
		key:    key,
		Logger: logger,
	}
}

func (d *Dictionary) query(text string) string {
	const GetUrl = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"

	text = url.PathEscape(text)
	return fmt.Sprintf(GetUrl, text, d.key)
}

func (d *Dictionary) Translate(text string) (data *DictionaryWord, err error) {
	query := d.query(text)
	d.Infof("query = %s, for text = %s", query, text)

	resp, err := http.Get(query)
	if err != nil {
		d.Errorf("error on http get %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d.Errorf("received status different than OK %s", resp.Status)
		return nil, fmt.Errorf("wrong status code %s", resp.Status)
	}

	data = &DictionaryWord{logger: d.Logger}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		d.Debugf("error decoding json to DictionarySingleWord %v", err)
		return nil, err
	}
	return
}

func (d *DictionaryWord) ShortDef() string {
	return strings.Join(d.Words[0].Shortdef, "\n")
}
func (d *DictionaryWord) Examples() string {
	s := make([]string, 0, len(d.Words[0].Suppl.Examples))
	for _, elem := range d.Words[0].Suppl.Examples {
		s = append(s, elem.T)
	}
	return strings.Join(s, "\n")
}

func (d *DictionaryWord) AudioFiles() []string {
	const AudioUrl = `https://media.merriam-webster.com/audio/prons/en/us/mp3/%s/%s.mp3`

	paths := make([]string, 0, len(d.Words[0].Hwi.Prs))
	for _, elem := range d.Words[0].Hwi.Prs {
		filename := elem.Sound.Audio
		dir := ""
		if strings.HasPrefix(filename, "bix") {
			dir = "bix"
		} else if strings.HasPrefix(filename, "gg") {
			dir = "gg"
		} else if unicode.IsNumber(rune(filename[0])) || unicode.IsPunct(rune(filename[0])) {
			dir = "number"
		} else {
			dir = string(filename[0])
		}
		paths = append(paths, fmt.Sprintf(AudioUrl, dir, filename))
	}
	return paths

}

type DictionarySingleWord struct {
	Meta struct {
		Id        string   `json:"id"`
		Uuid      string   `json:"uuid"`
		Sort      string   `json:"sort"`
		Src       string   `json:"src"`
		Section   string   `json:"section"`
		Stems     []string `json:"stems"`
		Offensive bool     `json:"offensive"`
	} `json:"meta"`
	Hwi struct {
		Hw  string `json:"hw"`
		Prs []struct {
			Mw    string `json:"mw"`
			Sound struct {
				Audio string `json:"audio"`
				Ref   string `json:"ref"`
				Stat  string `json:"stat"`
			} `json:"sound"`
		} `json:"prs"`
	} `json:"hwi"`
	Fl  string `json:"fl"`
	Def []struct {
		Sseq [][][]interface{} `json:"sseq"`
	} `json:"def"`
	Uros []struct {
		Ure string `json:"ure"`
		Fl  string `json:"fl"`
	} `json:"uros"`
	Et     [][]string `json:"et"`
	Date   string     `json:"date"`
	LdLink struct {
		LinkHw string `json:"link_hw"`
		LinkFl string `json:"link_fl"`
	} `json:"ld_link"`
	Suppl struct {
		Examples []struct {
			T string `json:"t"`
		} `json:"examples"`
		Ldq struct {
			Ldhw string `json:"ldhw"`
			Fl   string `json:"fl"`
			Def  []struct {
				Sls  []string          `json:"sls"`
				Sseq [][][]interface{} `json:"sseq"`
			} `json:"def"`
		} `json:"ldq"`
	} `json:"suppl"`
	Shortdef []string `json:"shortdef"`
}
