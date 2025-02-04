package parameters

type SettingsRepositoryMemory struct{}

func (s *SettingsRepositoryMemory) GetActiveProxySettings() *ProxyEngine {
	result := ProxyEngine{}
	// rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123", true}
	rmsMain := RMSConnectParameter{"rmsMain", "http://192.168.2.7:8080", "adm", "40bd001563085fc35165329ea1ff5c5ecbdbbeef", false}
	// rms 2
	// rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123", true}
	result.MainRms = rmsMain
	rmsFake := RMSConnectParameter{"fake", "https://proizvodstvo-yaponskii-karri.iiko.it", "iikoUser", "jwSCPwVvd3", true}
	// rmsFakeTwo := RMSConnectParameter{"https://rokets-dzhentl.iiko.it", "iikoUser", "y4Wce0qMgb", true}
	replaceOne := ReplacedItem{}
	replaceOne.PfakeRms = &rmsFake
	replaceOne.Path = "/resto/api/employees/schedule/byDepartment/3/"
	replaceOne.PathTo = "/resto/api/employees/schedule/byDepartment/9/"

	// replaceTwo := ReplacedItem{}
	// replaceTwo.fakeRms = &rmsFakeTwo
	// replaceTwo.path = "/resto/api/employees/schedule/byDepartment/3/"
	// repl := []ReplacedItem{replaceOne, replaceTwo}
	repl := []ReplacedItem{replaceOne}
	// repl := []ReplacedItem{}
	result.Replaced = repl
	return &result
}
