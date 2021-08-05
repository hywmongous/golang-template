package commands

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type DeleteIdentity struct {
	IdentityID values.IdentityID
}

func CreateDeleteIdentity(identityID values.IdentityID) DeleteIdentity {
	return DeleteIdentity{
		IdentityID: identityID,
	}
}

func (deletion DeleteIdentity) Apply(handler CommandHandler) error {
	return handler.VisitDeleteIdentity(deletion)
}
