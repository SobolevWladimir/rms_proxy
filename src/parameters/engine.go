package parameters

import "net/http"

type ProxyEngine struct {
	replaced []ReplacedItem      // список элементов на замену
	mainRms  RMSConnectParameter //  Основоной контент от куда берем данные
	logs     []LogItem
}

func (e *ProxyEngine) Handle(r *http.Request) (*http.Response, error) {
	// проверяем, если тут наш api.
	rep := e.getReplaceItem(r)
	if rep != nil {
		return rep.Handle(r)
	}
	return e.mainRms.Handle(r)
}

// Находим нужный элемент для замены
func (e *ProxyEngine) getReplaceItem(r *http.Request) *ReplacedItem {
	// обходим каждый элемн на соответсвие  ..и если что проверякм

	return nil
}
