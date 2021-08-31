package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type IdentityLoginRequest struct {
	IdentityID identity.IdentityID
	Password   string
}

type IdentityLoginResponse struct {
	Id string
}

type IdentityLoginUseCase func(request IdentityLoginRequest) (IdentityLoginResponse, error)
