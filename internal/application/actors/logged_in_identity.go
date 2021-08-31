package actors

import "github.com/hywmongous/example-service/internal/application/usecases"

type LoggedInIdentity struct {
	Logout usecases.IdentityLogoutUseCase
}

func CreateLoggedInIdentity() LoggedInIdentity {
	return LoggedInIdentity{}
}
