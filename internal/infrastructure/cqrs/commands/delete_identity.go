package commands

import identity "github.com/hywmongous/example-service/internal/domain/identity/values"

type DeleteIdentity struct {
	IdentityID identity.IdentityID
}

func CreateDeleteIdentity(identityID identity.IdentityID) DeleteIdentity {
	return DeleteIdentity{
		IdentityID: identityID,
	}
}

func (deletion DeleteIdentity) Apply(handler CommandHandler) error {
	return handler.VisitDeleteIdentity(deletion)
}
