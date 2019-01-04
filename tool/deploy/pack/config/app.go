package config

type App struct {
	Name string `json:"name"`
	Bin  Binary `json:"bin"`
	Src  Source `json:"src"`
}
