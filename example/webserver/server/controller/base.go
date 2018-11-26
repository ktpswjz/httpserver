package controller

import (
	"fmt"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/types"
)

type Base struct {
	types.Base

	Config  *config.Config
	DbToken memory.Token
}

func (s *Base) GetToken(token string) (*model.Token, error) {
	if token == "" {
		return nil, fmt.Errorf("empty token id")
	}

	return s.DbToken.Get(token)
}
