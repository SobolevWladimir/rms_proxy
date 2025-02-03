package localstore

import "github.com/goccy/go-json"

type ConfigReplacedItem struct {
	Path             string          `json:"path"`             // url который подменяем
	PathTo           string          `json:"pathTo"`           // url на который подменяем по умолчанию  path
	ReplaceByFakeRms bool            `json:"replaceByFakeRms"` // подменить запрос с помощью другой rms
	PfakeRmsID       json.RawMessage `json:"fakeRmsId"`        // само фейковое рмs
	PfakeContent     json.RawMessage `json:"fakeContent"`      //  ответ .. если не берем из  рмс
}
