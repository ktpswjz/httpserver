package model

import "time"

type Token struct {
	ID          string    `json:"id" note:"标识ID"`
	UserID      uint64    `json:"userId" note:"用户ID"`
	UserAccount string    `json:"userAccount" note:"用户账号"`
	UserName    string    `json:"userName" note:"用户姓名"`
	LoginIP     string    `json:"loginIp" note:"用户登陆IP"`
	LoginTime   time.Time `json:"loginTime" note:"登陆时间"`
	ActiveTime  time.Time `json:"activeTime" note:"最近激活时间"`
}

func (s *Token) CopyTo(target *Token) {
	if target == nil {
		return
	}

	target.ID = s.ID
	target.UserID = s.UserID
	target.UserAccount = s.UserAccount
	target.UserName = s.UserName
	target.LoginIP = s.LoginIP
	target.LoginTime = s.LoginTime
	target.ActiveTime = s.ActiveTime
}

type TokenFilter struct {
	Account string `json:"account"`
}
