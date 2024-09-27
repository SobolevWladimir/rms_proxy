package parameters

type SettingsRepositoryMemory struct {
}

func (s *SettingsRepositoryMemory) GetActiveProxySettings() *ProxyEngine {
	result := ProxyEngine{}
	rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123", true}
	// rmsFake := RMSConnectParameter{"https://rokets-co.iiko.it", "iikoUser", "ExJwoACOdD"}
	// rmsFakeTwo := RMSConnectParameter{"https://rokets-dzhentl.iiko.it", "iikoUser", "y4Wce0qMgb"}
	rmsFake := RMSConnectParameter{"https://rokets-co.iiko.it", "SystemIntegrationUser", "cae88cb3fc90e2e9224b67e48cf63642fc0b54d8", false}
	rmsFakeTwo := RMSConnectParameter{"https://rokets-dzhentl.iiko.it", "SystemIntegrationUser", "cae88cb3fc90e2e9224b67e48cf63642fc0b54d8", false}
	result.mainRms = rmsMain
	replaceOne := ReplacedItem{}
	replaceOne.fakeRms = &rmsFake
	replaceOne.path = "/resto/api/products/search"

	replaceTwo := ReplacedItem{}
	replaceTwo.fakeRms = &rmsFakeTwo
	replaceTwo.path = "/resto/api/v2/reports/olap"
	repl := []ReplacedItem{replaceOne, replaceTwo}
	result.replaced = repl
	return &result
}
