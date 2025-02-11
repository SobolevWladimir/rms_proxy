package localstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (c *ConfigStore) SaveRmsList(l ConfigRmsList) error {
	f, _ := os.Create(c.getFileForRmsSettings())
	defer f.Close()
	json, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = f.Write(json)
	return err
}

func (c *ConfigStore) GetRMSList() ConfigRmsList {
	if _, err := os.Stat(c.getFileForRmsSettings()); errors.Is(err, os.ErrNotExist) {
		return c.getDefaultRmsList()
	}
	content, err := os.ReadFile(c.getFileForRmsSettings())
	if err != nil {
		fmt.Println("Ошибка открытия файла rms.json ", err.Error())
		return c.getDefaultRmsList()
	}
	result := ConfigRmsList{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Ошибка чтения файла rms.json ", err.Error())
		return c.getDefaultRmsList()
	}
	return result
}

func (c *ConfigStore) getDefaultRmsList() ConfigRmsList {
	chain := "chain.rms_proxy.localhost"
	rms1 := "rms1.rms_proxy.localhost"
	rmsChain := ConfigRmsItem{
		ID:              "rmsChain",
		Name:            "Демо Chain",
		URL:             "http://192.168.0.155:9080",
		Login:           "adm",
		Password:        "123",
		NeedPassEncrupt: true,
		Domain:          &chain,
	}

	rmsOne := ConfigRmsItem{
		ID:              "rmsOne",
		Name:            "Демо RMS 1",
		URL:             "http://192.168.2.7:8080",
		Login:           "adm",
		Password:        "40bd001563085fc35165329ea1ff5c5ecbdbbeef",
		NeedPassEncrupt: false,
		Domain:          &rms1,
	}

	result := ConfigRmsList{List: []ConfigRmsItem{rmsChain, rmsOne}}
	return result
}
