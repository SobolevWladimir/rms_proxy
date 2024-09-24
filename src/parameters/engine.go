package parameters

import "net/http"

type ProxyEngine struct {
	replaced []ReplacedItem      // список элементов на замену
	mainRms  RMSConnectParameter //  Основоной контент от куда берем данные
	logs     []LogItem
}

func (e *ProxyEngine) Handle(r *http.Request) (string, error) {
	// проверяем, если тут наш api.
	rep := e.getReplaceItem(r)
	var result string
	var err error
	if rep != nil {
		result, err = rep.Handle(r)
	} else {
		result, err = e.mainRms.Handle(r)
	}
	// Создать лог.

	return result, err
}

// Находим нужный элемент для замены
func (e *ProxyEngine) getReplaceItem(r *http.Request) *ReplacedItem {
	// обходим каждый элемн на соответсвие  ..и если что проверякм

	return nil
}
