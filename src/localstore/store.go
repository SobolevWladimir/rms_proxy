package localstore

import "fmt"

type ConfigStore struct {
	Path string
}

func (c *ConfigStore) GetRMSList() ConfigRmsList {
	return c.getDefaultRmsList()
}

func (c *ConfigStore) SaveRmsList(l ConfigRmsList) {
	fmt.Println("save list", l.MainRms)
}

func (c *ConfigStore) getDefaultRmsList() ConfigRmsList {
	rmsChain := ConfigRmsItem{
		ID:              "rmsChain",
		Name:            "Демо Chain",
		URL:             "http://192.168.0.155:9080",
		Login:           "adm",
		Password:        "123",
		NeedPassEncrupt: true,
	}

	rmsOne := ConfigRmsItem{
		ID:              "rmsOne",
		Name:            "Демо RMS 1",
		URL:             "http://192.168.2.7:8080",
		Login:           "adm",
		Password:        "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
		NeedPassEncrupt: false,
	}

	result := ConfigRmsList{List: []ConfigRmsItem{rmsChain, rmsOne}, MainRms: "rmsChain"}
	return result
}

