package config

type Site struct {
	App   SiteApp   `json:"app"`
	Doc   SiteDoc   `json:"doc"`
	Admin SiteAdmin `json:"admin"`
}
