package identity

type SessionContext struct {
	id             SessionContextID
	issuedAt       Time
	csrf           Csrf
	accessTokenId  AccessTokenID
	refreshTokenId RefreshTokenID
}

func CreateSessionContext() SessionContext {
	return SessionContext{
		id:             GenerateSessionContextID(),
		issuedAt:       Now(),
		csrf:           GenerateCsrf(),
		accessTokenId:  GenerateAccessTokenID(),
		refreshTokenId: GenerateRefreshTokenID(),
	}
}

func RecreateSessionContext(
	id SessionContextID,
	issuedAt Time,
	csrf Csrf,
	accessTokenId AccessTokenID,
	refreshTokenId RefreshTokenID,
) SessionContext {
	return SessionContext{
		id:             id,
		issuedAt:       issuedAt,
		csrf:           csrf,
		accessTokenId:  accessTokenId,
		refreshTokenId: refreshTokenId,
	}
}

func (context SessionContext) GetId() Csrf {
	return context.csrf
}

func (context SessionContext) GetIssuedAt() Time {
	return context.issuedAt
}

func (context SessionContext) GetCsrf() Csrf {
	return context.csrf
}

func (context SessionContext) GetAccessTokenId() AccessTokenID {
	return context.accessTokenId
}

func (context SessionContext) GetRefreshTokenId() RefreshTokenID {
	return context.refreshTokenId
}
