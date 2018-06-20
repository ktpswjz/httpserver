package config

type SiteAdmin struct {
	Enable bool `json:"enable"`
	Root string `json:"root"`
	Api SiteAdminApi `json:"api"`
	Users []SiteAdminUser `json:"users"`
	Ldap SiteAdminLdap `json:"ldap"`
}
