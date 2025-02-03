package localstore

import (
	"path"
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
