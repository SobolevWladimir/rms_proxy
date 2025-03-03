package parameters

import (
	"io"
	"net/http"
	"strings"
)

type ReplacedItem struct {
	Path              string               `json:"path"`                // url который подменяем
	Content           string               `json:"content"`             // Условие для совмещения контента
	QueryKeys         map[string]string    `json:"query_keys"`          // Условия  ключи в запросе
	IsContentContains bool                 `json:"is_content_contains"` // Условия Искать часть контента
	PathTo            string               `json:"pathTo"`              // url на который подменяем по умолчанию  path
	ReplaceByFakeRms  bool                 `json:"replaceByFakeRms"`    // подменить запрос с помощью другой rms
	PfakeRms          *RMSConnectParameter `json:"fakeRms"`             // само фейковое рмs
	PfakeContent      string               `json:"fakeContent"`         //  ответ .. если не берем из  рмс
}

func (rm *ReplacedItem) Handle(r *http.Request, log *LogItem) (*http.Response, error) {
	if len(rm.PathTo) > 0 {
		r.URL.Path = rm.PathTo
	}
	if rm.ReplaceByFakeRms {
		return rm.PfakeRms.Handle(r, log)
	}
	header := make(http.Header)

	header["test"] = []string{"rms_proxy"}

	result := http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(rm.PfakeContent)),
		Header:     header,
	}

	return &result, nil
}

func (rm *ReplacedItem) IsSuitable(r *http.Request) bool {
	if r.URL.Path != rm.Path {
		return false
	}
	if len(rm.Content) > 0 {
		b, _ := io.ReadAll(r.Body)
		body := string(b)
		r.Body = io.NopCloser(strings.NewReader(body))

		if rm.IsContentContains {
			if !strings.Contains(body, rm.Content) {
				return false
			}
		} else {
			if rm.Content != body {
				return false
			}
		}
	}
	query := r.URL.Query()
	for val, k := range rm.QueryKeys {
		if query.Has(k) {
			qVal := query.Get(k)
			if qVal != val {
				return false
			}
		}
	}

	return true
}
