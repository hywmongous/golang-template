package domain

import (
	"time"

	"github.com/hywmongous/example-service/pkg/guid"
)

type Session struct {
	Id           string
	Csrf         string
	Revoked      bool
	AccessToken  Token
	RefreshToken Token
}

type Token struct {
	Id              string
	IssuedAt        int64
	InitialTimeout  int64
	AbsoluteTimeout int64
}

const (
	AccessTokenInitialTimeoutDuration   = 0
	AccessTokenAbsoluteTimeoutDuration  = 30
	RefreshTokenInitialTimeoutDuration  = 15
	RefreshTokenAbsoluteTimeoutDuration = 30
)

func SessionFactory() Session {
	return Session{
		Id:           guid.New().String(),
		Csrf:         guid.New().String(),
		Revoked:      false,
		AccessToken:  accessTokenFactory(),
		RefreshToken: refreshTokenFactory(),
	}
}

func (session *Session) Refresh() {
	session.Csrf = guid.New().String()
	session.AccessToken = accessTokenFactory()
	session.RefreshToken = refreshTokenFactory()
}

func accessTokenFactory() Token {
	// Access tokens can be used immediately and expires after 30 minutes
	now := time.Now()
	return Token{
		Id:              guid.New().String(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(AccessTokenInitialTimeoutDuration * time.Minute).Unix(),
		AbsoluteTimeout: now.Add(AccessTokenAbsoluteTimeoutDuration * time.Minute).Unix(),
	}
}

func refreshTokenFactory() Token {
	// Refresh tokens can be used after 15 minutes and expires after 30
	now := time.Now()
	return Token{
		Id:              guid.New().String(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(RefreshTokenInitialTimeoutDuration * time.Minute).Unix(),
		AbsoluteTimeout: now.Add(RefreshTokenAbsoluteTimeoutDuration * time.Minute).Unix(),
	}
}
