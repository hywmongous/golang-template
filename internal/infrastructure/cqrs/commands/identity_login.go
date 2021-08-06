package commands

import identity "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityLogin struct {
	IdentityID identity.IdentityID
	Password   string
}

func CreateIdentityLogin(
	identityId identity.IdentityID,
	password string,
) IdentityLogin {
	return IdentityLogin{
		IdentityID: identityId,
		Password:   password,
	}
}

func (login IdentityLogin) Apply(handler CommandHandler) error {
	return handler.VisitIdentityLogin(login)
}
