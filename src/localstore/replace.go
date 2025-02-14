package localstore

import (
	"rms_proxy/v2/src/parameters"

	"github.com/goccy/go-json"
)

type ConfigReplacedItem struct {
	Path              string            `json:"path"`                // url который подменяем
	QueryKeys         map[string]string `json:"query_keys"`          // Условия  ключи в запросе
	IsContentContains bool              `json:"is_content_contains"` // Условия Искать часть контента
	Content           string            `json:"content"`             // Условие  Контент который должен быть
	PathTo            string            `json:"pathTo"`              // url на который подменяем по умолчанию  path
	ReplaceByFakeRms  bool              `json:"replaceByFakeRms"`    // подменить запрос с помощью другой rms
	PfakeRmsID        json.RawMessage   `json:"fakeRmsId"`           // само фейковое рмs
	PfakeContent      json.RawMessage   `json:"fakeContent"`         //  ответ .. если не берем из  рмс
}

func (c *ConfigReplacedItem) ToReplaceItem(rmsList map[string]*parameters.RMSConnectParameter) parameters.ReplacedItem {
	rms, ok := rmsList[string(c.PfakeRmsID)]
	if !ok {
		rms = nil
	}
	var content string
	var text string
	err := json.Unmarshal([]byte(c.PfakeContent), &text)
	if err != nil {
		text = string(c.PfakeContent)
	}
	err = json.Unmarshal([]byte(c.Content), &content)
	if err != nil {
		content = c.Content
	}
	result := parameters.ReplacedItem{
		Path:              c.Path,
		QueryKeys:         c.QueryKeys,
		IsContentContains: c.IsContentContains,
		Content:           content,
		PathTo:            c.PathTo,
		ReplaceByFakeRms:  c.ReplaceByFakeRms,
		PfakeRms:          rms,
		PfakeContent:      text,
	}
	return result
}
