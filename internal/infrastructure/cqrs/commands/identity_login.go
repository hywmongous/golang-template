package commands

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityLogin struct {
	IdentityID values.IdentityID
	Password   string
}

func CreateIdentityLogin(
	identityId values.IdentityID,
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
