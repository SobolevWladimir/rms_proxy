package localstore

type ConfigReplacedItem struct {
	Path             string `json:"path"`             // url который подменяем
	PathTo           string `json:"pathTo"`           // url на который подменяем по умолчанию  path
	ReplaceByFakeRms bool   `json:"replaceByFakeRms"` // подменить запрос с помощью другой rms
	PfakeRmsID       string `json:"fakeRmsId"`        // само фейковое рмs
	PfakeContent     string `json:"fakeContent"`      //  ответ .. если не берем из  рмс
}





