package localstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (c *ConfigStore) GetProxyItems() []ConfigReplacedItem {
	result := []ConfigReplacedItem{}
	if _, err := os.Stat(c.getFileForProxyItems()); errors.Is(err, os.ErrNotExist) {
		return result
	}
	content, err := os.ReadFile(c.getFileForProxyItems())
	if err != nil {
		fmt.Println("Ошибка открытия файла proxy.json ", err.Error())
		return result
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Ошибка чтения файла proxy.json ", err.Error())
		return result
	}
	return result
}

func (c *ConfigStore) SaveProxyItems(items []ConfigReplacedItem) error {
	f, _ := os.Create(c.getFileForProxyItems())
	defer f.Close()
	json, err := json.Marshal(items)
	if err != nil {
		return err
	}
	_, err = f.Write(json)
	return err
}
