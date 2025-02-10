package parameters

import (
	"fmt"
	"net/http"
)

type ProxyEngine struct {
	Replaced []ReplacedItem      // список элементов на замену
	MainRms  RMSConnectParameter //  Основоной контент от куда берем данные
	Port     string
}

func (e *ProxyEngine) Handle(r *http.Request) (*http.Response, LogItem) {
	log := LogItem{
		ClientRequest: CreateHTTPRequest(r),
	}
	fmt.Println("HOST: ", r.Host)
	// проверяем, если тут наш api.
	fmt.Println("host -------------------------------------")
	fmt.Println(r.Host);
	rep := e.getReplaceItem(r)
	var res *http.Response
	var err error
	log.MainRMS = e.MainRms
	if rep != nil {
		res, err = rep.Handle(r, &log)
		log.IsProxy = true
		log.ProxyTo = rep
	} else {
		res, err = e.MainRms.Handle(r, &log)
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
	for _, val := range e.Replaced {
		if val.IsSuitable(r) {
			return &val
		}
	}

	return nil
}
