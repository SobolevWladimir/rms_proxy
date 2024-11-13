package parameters

type SettingsRepositoryMemory struct {
}

func (s *SettingsRepositoryMemory) GetActiveProxySettings() *ProxyEngine {
	result := ProxyEngine{}
	// rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123", true}
	// rmsMain := RMSConnectParameter{"http://192.168.2.7:8080", "adm", "123", true}
	// rms 2
	rmsMain := RMSConnectParameter{"http://192.168.0.155:8080", "adm", "123", true}
	result.mainRms = rmsMain
	rmsFake := RMSConnectParameter{"https://proizvodstvo-yaponskii-karri.iiko.it", "iikoUser", "jwSCPwVvd3", true}
	// rmsFakeTwo := RMSConnectParameter{"https://rokets-dzhentl.iiko.it", "iikoUser", "y4Wce0qMgb", true}
	replaceOne := ReplacedItem{}
	replaceOne.fakeRms = &rmsFake
	replaceOne.path = "/resto/api/employees/schedule/byDepartment/3/"
	replaceOne.pathTo = "/resto/api/employees/schedule/byDepartment/9/"

	// replaceTwo := ReplacedItem{}
	// replaceTwo.fakeRms = &rmsFakeTwo
	// replaceTwo.path = "/resto/api/employees/schedule/byDepartment/3/"
	// repl := []ReplacedItem{replaceOne, replaceTwo}
	repl := []ReplacedItem{replaceOne}
	// repl := []ReplacedItem{}
	result.replaced = repl
	return &result
}
