package parameters

type SettingsRepositoryMemory struct {
}

func (s *SettingsRepositoryMemory) GetActiveProxySettings() *ProxyEngine {
	result := ProxyEngine{}
	rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123", true}
	result.mainRms = rmsMain
	// rmsFake := RMSConnectParameter{"https://rokets-co.iiko.it", "iikoUser", "ExJwoACOdD", true}
	// rmsFakeTwo := RMSConnectParameter{"https://rokets-dzhentl.iiko.it", "iikoUser", "y4Wce0qMgb", true}
	// replaceOne := ReplacedItem{}
	// replaceOne.fakeRms = &rmsFake
	// replaceOne.path = "/resto/api/products/search"

	// replaceTwo := ReplacedItem{}
	// replaceTwo.fakeRms = &rmsFakeTwo
	// replaceTwo.path = "/resto/api/v2/reports/olap"
	// repl := []ReplacedItem{replaceOne, replaceTwo}
	repl := []ReplacedItem{}
	result.replaced = repl
	return &result
}
