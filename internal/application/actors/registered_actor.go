package actors

import (
	"github.com/hywmongous/example-service/internal/application/usecases"
	"github.com/hywmongous/example-service/internal/domain/identity/values"
)

type Registered struct {
	RegisteredIdentity values.IdentityID

	unregistered Unregistered
	login        usecases.Login
	logout       usecases.Logout
	delete       usecases.Delete
}

func CreateRegistered(
	registeredIdentity values.IdentityID,
	unregistered Unregistered,
	login usecases.Login,
	logout usecases.Logout,
	delete usecases.Delete,
) (Registered, error) {
	return Registered{
		RegisteredIdentity: registeredIdentity,
		unregistered:       unregistered,
		login:              login,
		logout:             logout,
		delete:             delete,
	}, nil
}
