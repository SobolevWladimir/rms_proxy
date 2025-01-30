package parameters

import (
	"net/http"
)

type ProxyEngine struct {
	replaced []ReplacedItem      // список элементов на замену
	mainRms  RMSConnectParameter //  Основоной контент от куда берем данные
}

func (e *ProxyEngine) Handle(r *http.Request) (*http.Response, LogItem) {
	log := LogItem{
		ClientRequest: CreateHTTPRequest(r),
	}
	// проверяем, если тут наш api.
	rep := e.getReplaceItem(r)
	var res *http.Response
	var err error
	log.MainRMS = e.mainRms
	if rep != nil {
		res, err = rep.Handle(r)
		log.IsProxy = true
		log.ProxyTo = rep
	} else {
		res, err = e.mainRms.Handle(r)
		log.IsProxy = false
		log.ProxyTo = nil
	}
	log.IsErrorResponse = err != nil
	log.ClientResponse = CreateHTTPResponse(res)
	if err != nil {
		log.ErrorResponse = err.Error()
	}

	return res, log
}

// Находим нужный элемент для замены
func (e *ProxyEngine) getReplaceItem(r *http.Request) *ReplacedItem {
	// обходим каждый элемн на соответсвие  ..и если что проверякм
	for _, val := range e.replaced {
		if val.IsSuitable(r) {
			return &val
		}
	}

	return nil
}
