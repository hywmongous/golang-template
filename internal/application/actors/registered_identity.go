package actors

import "github.com/hywmongous/example-service/internal/application/usecases"

type RegisteredIdentity struct {
	Login  usecases.IdentityLoginUseCase
	Delete usecases.DeleteIdentityUseCase
}

func CreateRegisteredIdentity() RegisteredIdentity {
	return RegisteredIdentity{}
}
