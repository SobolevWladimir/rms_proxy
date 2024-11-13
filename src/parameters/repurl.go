package parameters

import (
	"fmt"
	"net/http"
)

type ReplacedItem struct {
	path             string               // url который подменяем
	pathTo           string               // url на который подменяем по умолчанию  path
	replaceByFakeRms bool                 // подменить запрос с помощью другой rms
	fakeRms          *RMSConnectParameter // само фейковое рмs
	fakeContent      string               //  ответ .. если не берем из  рмс
}

func (rm *ReplacedItem) Handle(r *http.Request) (*http.Response, error) {
	fmt.Println("\n Handle: ", r.URL.Path)
	if(len(rm.pathTo)>0){
		r.URL.Path = rm.pathTo
	}
	return rm.fakeRms.Handle(r)
}

func (rm *ReplacedItem) IsSuitable(r *http.Request) bool {
	if r.URL.Path == rm.path {
		return true
	}
	return false
}
