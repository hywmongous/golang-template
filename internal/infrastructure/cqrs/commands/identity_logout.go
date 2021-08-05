package commands

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityLogout struct {
	IdentityID values.IdentityID
	SessionID  values.SessionID
}

func CreateIdentityLogou(
	identityId values.IdentityID,
	sessionID values.SessionID,
) IdentityLogout {
	return IdentityLogout{
		IdentityID: identityId,
		SessionID:  sessionID,
	}
}

func (logout IdentityLogout) Apply(handler CommandHandler) error {
	return handler.VisitIdentityLogout(logout)
}
