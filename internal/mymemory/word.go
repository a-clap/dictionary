package mymemory

func (w *Word) Translated() string {
	return w.ResponseData.TranslatedText
}

type Alternative struct {
	text, translation string
}

func (w *Word) Alternatives() []Alternative {
	alt := make([]Alternative, 0, len(w.Matches))
	for _, elem := range w.Matches {
		alt = append(alt, Alternative{
			text:        elem.Segment,
			translation: elem.Translation,
		})
	}
	return alt
}

type Word struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
		//Match          float64 `json:"match"`
	} `json:"responseData"`
	//QuotaFinished   bool        `json:"quotaFinished"`
	//MtLangSupported interface{} `json:"mtLangSupported"`
	//ResponseDetails string      `json:"responseDetails"`
	//ResponseStatus  int         `json:"responseStatus"`
	//ResponderId     string      `json:"responderId"`
	//ExceptionCode   interface{} `json:"exception_code"`
	Matches []struct {
		//Id             interface{} `json:"id"`
		Segment     string `json:"segment"`
		Translation string `json:"translation"`
		//Source         string      `json:"source"`
		//Target         string      `json:"target"`
		//Quality        string      `json:"quality"`
		//Reference      interface{} `json:"reference"`
		//UsageCount     int         `json:"usage-count"`
		//Subject        string      `json:"subject"`
		//CreatedBy      string      `json:"created-by"`
		//LastUpdatedBy  string      `json:"last-updated-by"`
		//CreateDate     string      `json:"create-date"`
		//LastUpdateDate string      `json:"last-update-date"`
		//Match          float64     `json:"match"`
	} `json:"matches"`
}
