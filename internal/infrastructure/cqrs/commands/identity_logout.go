package commands

import identity "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityLogout struct {
	IdentityID identity.IdentityID
	SessionID  identity.SessionID
}

func CreateIdentityLogou(
	identityId identity.IdentityID,
	sessionID identity.SessionID,
) IdentityLogout {
	return IdentityLogout{
		IdentityID: identityId,
		SessionID:  sessionID,
	}
}

func (logout IdentityLogout) Apply(handler CommandHandler) error {
	return handler.VisitIdentityLogout(logout)
}
