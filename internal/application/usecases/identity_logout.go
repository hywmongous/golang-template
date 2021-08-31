package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type IdentityLogoutRequest struct {
	IdentityID identity.IdentityID
	Password   string
}

type IdentityLogoutResponse struct {
}

type IdentityLogoutUseCase func(request IdentityLogoutRequest) (IdentityLoginResponse, error)
