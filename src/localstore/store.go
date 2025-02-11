package localstore

import (
	"path"

	"rms_proxy/v2/src/parameters"
)

type ConfigStore struct {
	Path string
}

func (c *ConfigStore) getFileForRmsSettings() string {
	pwd := path.Join(c.Path, "rms.json")
	return pwd
}

func (c *ConfigStore) getFileForProxyItems() string {
	pwd := path.Join(c.Path, "proxy.json")
	return pwd
}

func (c *ConfigStore) GetActiveProxySettings() *parameters.ProxyEngine {
	result := parameters.ProxyEngine{}
	rmsList := c.GetRMSList()
	result.RmsList = rmsList.ToParameter()
	result.Port = ":9091"
	rmsMap := c.getRms()

	items := c.GetProxyItems()
	for _, item := range items {
		rep := item.ToReplaceItem(rmsMap)
		result.Replaced = append(result.Replaced, rep)
	}
	return &result
}

func (c *ConfigStore) getRms() map[string]*parameters.RMSConnectParameter {
	result := make(map[string]*parameters.RMSConnectParameter)
	rmsList := c.GetRMSList()
	for _, con := range rmsList.List {
		item := con.ToParameter()
		result[con.ID] = &item
	}

	return result
}
