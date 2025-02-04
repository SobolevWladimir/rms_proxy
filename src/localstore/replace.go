package localstore

import (
	"rms_proxy/v2/src/parameters"

	"github.com/goccy/go-json"
)

type ConfigReplacedItem struct {
	Path             string          `json:"path"`             // url который подменяем
	PathTo           string          `json:"pathTo"`           // url на который подменяем по умолчанию  path
	ReplaceByFakeRms bool            `json:"replaceByFakeRms"` // подменить запрос с помощью другой rms
	PfakeRmsID       json.RawMessage `json:"fakeRmsId"`        // само фейковое рмs
	PfakeContent     json.RawMessage `json:"fakeContent"`      //  ответ .. если не берем из  рмс
}

func (c *ConfigReplacedItem) ToReplaceItem(rmsList map[string]*parameters.RMSConnectParameter) parameters.ReplacedItem {
	rms, ok := rmsList[string(c.PfakeRmsID)]
	if !ok {
		rms = nil
	}

	result := parameters.ReplacedItem{
		Path:             c.Path,
		PathTo:           c.PathTo,
		ReplaceByFakeRms: c.ReplaceByFakeRms,
		PfakeRms:         rms,
		PfakeContent:     string(c.PfakeContent),
	}
	return result
}
