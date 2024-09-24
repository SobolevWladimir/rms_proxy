package parameters

type SettingsRepositoryMemory struct {
}

func (s *SettingsRepositoryMemory) GetActiveProxySettings() *ProxyEngine {
	result := ProxyEngine{}
	rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123"}
	rmsFake := RMSConnectParameter{"https://hizhina-pavlodar.iiko.it", "iikoUser", "NsoXTO9HuP"}
	result.mainRms = rmsMain
	replaceOne := ReplacedItem{}
	replaceOne.fakeRms = &rmsFake
	replaceOne.path = "/resto/api/v3/EntitiesService.getEntitiesUpdateddfsd"
	repl := []ReplacedItem{replaceOne}
	result.replaced = repl
	return &result
}
