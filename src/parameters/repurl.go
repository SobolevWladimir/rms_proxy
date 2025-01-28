package parameters

import (
	"fmt"
	"net/http"
)

type ReplacedItem struct {
	Path             string             `json:"path"`  // url который подменяем
	PathTo           string              `json:"pathTo"` // url на который подменяем по умолчанию  path
	ReplaceByFakeRms bool                   `json:"replaceByFakeRms"` // подменить запрос с помощью другой rms
	PfakeRms          *RMSConnectParameter `json:"fakeRms"` // само фейковое рмs
	PfakeContent      string             `json:"fakeContent"`  //  ответ .. если не берем из  рмс
}

func (rm *ReplacedItem) Handle(r *http.Request) (*http.Response, error) {
	fmt.Println("\n Handle: ", r.URL.Path)
	if(len(rm.PathTo)>0){
		r.URL.Path = rm.PathTo
	}
	return rm.PfakeRms.Handle(r)
}

func (rm *ReplacedItem) IsSuitable(r *http.Request) bool {
	if r.URL.Path == rm.Path {
		return true
	}
	return false
}
