package commands

import identity "github.com/hywmongous/example-service/internal/domain/identity/values"

type RegisterIdentity struct {
	IdentityID   identity.IdentityID
	EmailAddress string
	Password     string
}

func CreateRegisterIdentity(
	identityId identity.IdentityID,
	email string,
	password string,
) RegisterIdentity {
	return RegisterIdentity{
		IdentityID:   identityId,
		EmailAddress: email,
		Password:     password,
	}
}

func (registration RegisterIdentity) Apply(handler CommandHandler) error {
	return handler.VisitRegisterIdentity(registration)
}
