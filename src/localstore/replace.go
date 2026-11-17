package localstore

import (
	"encoding/json"

	"rms_proxy/v2/src/parameters"
)

type ConfigReplacedItem struct {
	Path              string            `json:"path"`                // url который подменяем
	QueryKeys         map[string]string `json:"query_keys"`          // Условия  ключи в запросе
	IsContentContains bool              `json:"is_content_contains"` // Условия Искать часть контента
	Content           string            `json:"content"`             // Условие  Контент который должен быть
	PathTo            string            `json:"pathTo"`              // url на который подменяем по умолчанию  path
	ReplaceByFakeRms  bool              `json:"replaceByFakeRms"`    // подменить запрос с помощью другой rms
	PfakeRmsID        json.RawMessage   `json:"fakeRmsId"`           // само фейковое рмs
	HasContent        bool              `json:"hasContent"`          // есть ли файл с подменным контентом
}

type ConfigReplacedItemInput struct {
	ConfigReplacedItem
	FakeContent string `json:"fakeContent"`
}

func rawJSONString(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return string(raw)
	}
	return s
}

func (c *ConfigReplacedItem) ToReplaceItem(rmsList map[string]*parameters.RMSConnectParameter, contentFile string) parameters.ReplacedItem {
	rms, ok := rmsList[rawJSONString(c.PfakeRmsID)]
	if !ok {
		rms = nil
	}
	result := parameters.ReplacedItem{
		Path:              c.Path,
		QueryKeys:         c.QueryKeys,
		IsContentContains: c.IsContentContains,
		Content:           rawJSONString(json.RawMessage(c.Content)),
		PathTo:            c.PathTo,
		ReplaceByFakeRms:  c.ReplaceByFakeRms,
		PfakeRms:          rms,
		ContentFile:       contentFile,
	}
	return result
}
