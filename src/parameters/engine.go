package parameters

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ProxyEngine struct {
	Replaced []ReplacedItem        // список элементов на замену
	RmsList  []RMSConnectParameter //  Основоной контент от куда берем данные
	Port     string
}

func (e *ProxyEngine) Handle(r *http.Request) (*http.Response, LogItem) {
	log := LogItem{
		ClientRequest: CreateHTTPRequest(r),
	}
	fmt.Println("HOST: ", r.Host)
	// проверяем, если тут наш api.
	fmt.Println(r.Host)
	rep := e.getReplaceItem(r)
	var res *http.Response
	var err error
	mainRms := e.getMainRms(r.Host)
	log.MainRMS = mainRms
	if mainRms == nil {
		notRms := RMSConnectParameter{Name: "not found"}
		log.MainRMS = &notRms
	}
	if rep != nil {
		res, err = rep.Handle(r, &log)
		log.IsProxy = true
		log.ProxyTo = rep
	} else {
		if mainRms != nil {
			res, err = mainRms.Handle(r, &log)
		} else {
			res = e.createResponseNonFoundRms()
			err = errors.New("Не смог найти  rms")
		}
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

func (e *ProxyEngine) createResponseNonFoundRms() *http.Response {
	body := string("Не найден рмс по домену")
	result := http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	return &result
}

func (e *ProxyEngine) getMainRms(host string) *RMSConnectParameter {
	fmt.Println("find")
	hostParts := strings.Split(host, ".")
	if len(hostParts) == 0 {
		fmt.Println("не указан host в заголовке")
		return nil
	}
	for _, rms := range e.RmsList {
		fmt.Println("+++++++++++++++++++++++++++++++")
		fmt.Println(rms.Domain, hostParts[0])
		if rms.Domain == hostParts[0] {
			return &rms
		}
	}

	return nil
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
