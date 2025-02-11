package localstore

import "rms_proxy/v2/src/parameters"

type ConfigRmsList struct {
	List []ConfigRmsItem `json:"list"`
}

type ConfigRmsItem struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	URL             string  `json:"url"`
	Login           string  `json:"login"`
	Password        string  `json:"password"`
	NeedPassEncrupt bool    `json:"needPassEncrupt"`
	Domain          *string `json:"domain"`
}

func (c *ConfigRmsItem) ToParameter() parameters.RMSConnectParameter {
	result := parameters.RMSConnectParameter{}
	result.Name = c.Name
	result.URL = c.URL
	result.Login = c.Login
	result.Password = c.Password
	result.NeedPassEncrupt = c.NeedPassEncrupt
	result.Domain = *c.Domain
	return result
}

func (c *ConfigRmsList) ToParameter() []parameters.RMSConnectParameter {
	result := []parameters.RMSConnectParameter{}
	for _, rms := range c.List {
		result = append(result, rms.ToParameter())
	}
	return result
}
