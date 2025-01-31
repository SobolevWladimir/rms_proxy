package localstore

type ConfigRmsList struct {
	List    []ConfigRmsItem `json:"list"`
	MainRms string          `json:"mainRmsId"` // RMS по умолчанию
}

type ConfigRmsItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Login           string `json:"login"`
	Password        string `json:"password"`
	NeedPassEncrupt bool   `json:"needPassEncrupt"`
}

