package localstore

import (
	"path"

	"rms_proxy/v2/src/parameters"
)

type ConfigStore struct {
	Path string
}
type ConfigProxyEngine struct {
	MainRms     ConfigRmsItem
	configStore *ConfigStore
}

func (c *ConfigStore) getFileForRmsSettings() string {
	pwd := path.Join(c.Path, "rms.json")
	return pwd
}

func (c *ConfigStore) getFileForProxyItems() string {
	pwd := path.Join(c.Path, "proxy.json")
	return pwd
}

func (c *ConfigProxyEngine) GetActiveProxySettings() *parameters.ProxyEngine {
	result := parameters.ProxyEngine{}
	rmsMap := c.getRms()
	// TODO: тут остановился
	// Надо еще перезагрузку сделать 

	return &result
}

func (c *ConfigProxyEngine) getRms() map[string]*parameters.RMSConnectParameter {
	result := make(map[string]*parameters.RMSConnectParameter)
	rmsList := c.configStore.GetRMSList()
	for _, con := range rmsList.List {
		item := con.ToParameter()
		result[con.ID] = &item
	}

	return result
}
