package authentication

import (
	"errors"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"strings"
)

type Admin struct {
	Config *config.SiteAdmin
}

func (s *Admin) Authenticate(account, password string) error {
	if s.Config == nil {
		return errors.New("internal error: nil config for admin site")
	}

	var user *config.SiteAdminUser = nil
	userCount := len(s.Config.Users)
	for index := 0; index < userCount; index++ {
		if account == strings.ToLower(s.Config.Users[index].Account) {
			user = &s.Config.Users[index]
			break
		}
	}

	if user != nil {
		if user.Password != password {
			return errors.New("invalid password")
		} else {
			return nil
		}
	}

	if !s.Config.Ldap.Enable {
		return errors.New("account not exist")
	}

	ldap := &Ldap{
		Host: s.Config.Ldap.Host,
		Port: s.Config.Ldap.Port,
		Base: s.Config.Ldap.Base,
	}

	return ldap.Authenticate(account, password)
}
