package aggregate

import (
	"github.com/hywmongous/example-service/internal/domain/identity/values"
)

type SessionContext struct {
	id             values.SessionContextID
	issuedAt       values.Time
	csrf           values.Csrf
	accessTokenId  values.AccessTokenID
	refreshTokenId values.RefreshTokenID
}

func CreateSessionContext() SessionContext {
	return SessionContext{
		id:             values.GenerateSessionContextID(),
		issuedAt:       values.Now(),
		csrf:           values.GenerateCsrf(),
		accessTokenId:  values.GenerateAccessTokenID(),
		refreshTokenId: values.GenerateRefreshTokenID(),
	}
}

func RecreateSessionContext(
	id values.SessionContextID,
	issuedAt values.Time,
	csrf values.Csrf,
	accessTokenId values.AccessTokenID,
	refreshTokenId values.RefreshTokenID,
) SessionContext {
	return SessionContext{
		id:             id,
		issuedAt:       issuedAt,
		csrf:           csrf,
		accessTokenId:  accessTokenId,
		refreshTokenId: refreshTokenId,
	}
}

func (context SessionContext) GetId() values.Csrf {
	return context.csrf
}

func (context SessionContext) GetIssuedAt() values.Time {
	return context.issuedAt
}

func (context SessionContext) GetCsrf() values.Csrf {
	return context.csrf
}

func (context SessionContext) GetAccessTokenId() values.AccessTokenID {
	return context.accessTokenId
}

func (context SessionContext) GetRefreshTokenId() values.RefreshTokenID {
	return context.refreshTokenId
}
