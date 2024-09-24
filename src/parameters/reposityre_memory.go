package parameters


type SettingsRepositoryMemory struct {

}

func(s *SettingsRepositoryMemory) GetActiveProxySettings()*ProxyEngine {
	result :=ProxyEngine{}
	rmsMain := RMSConnectParameter{"http://192.168.0.155:9080", "adm", "123"}
	result.mainRms = rmsMain
	return &result
}
