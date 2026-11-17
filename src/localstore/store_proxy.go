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

func (c *ConfigStore) SaveProxyItems(items []ConfigReplacedItemInput) error {
	meta := make([]ConfigReplacedItem, 0, len(items))
	for _, item := range items {
		it := item.ConfigReplacedItem
		if len(item.FakeContent) > 0 {
			it.HasContent = true
			if err := c.saveProxyContent(it.Path, item.FakeContent); err != nil {
				return err
			}
		} else {
			it.HasContent = false
			_ = c.deleteProxyContent(it.Path)
		}
		meta = append(meta, it)
	}

	f, err := os.Create(c.getFileForProxyItems())
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (c *ConfigStore) GetProxyContent(urlPath string) (string, error) {
	content, err := os.ReadFile(c.getContentFilePath(urlPath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (c *ConfigStore) saveProxyContent(urlPath, content string) error {
	if err := os.MkdirAll(c.getContentDir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(c.getContentFilePath(urlPath), []byte(content), 0o644)
}

func (c *ConfigStore) deleteProxyContent(urlPath string) error {
	return os.Remove(c.getContentFilePath(urlPath))
}
