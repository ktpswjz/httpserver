package config

type SiteAdminLdap struct {
	Enable 	bool 	`json:"enable"`
	Host 	string	`json:"host"`
	Port	int		`json:"port"`
	Base 	string	`json:"base"`
}
