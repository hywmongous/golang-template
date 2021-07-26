package domain

import (
	"time"

	"github.com/hywmongous/example-service/pkg/guid"
)

type SessionContext struct {
	Id             string
	IssuedAt       int64
	Csrf           string
	AccessTokenId  string
	RefreshTokenId string
}

func SessionContextFactory() SessionContext {
	return SessionContext{
		Id:             guid.New().String(),
		IssuedAt:       time.Now().Unix(),
		Csrf:           guid.New().String(),
		AccessTokenId:  guid.New().String(),
		RefreshTokenId: guid.New().String(),
	}
}
