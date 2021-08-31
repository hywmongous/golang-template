package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type RegisterIdentityRequest struct {
	Email    identity.Email
	Password identity.Password
}

type RegisterIdentityResponse struct {
	IdentityID identity.IdentityID
}

type RegisterIdentityUseCase func(request RegisterIdentityRequest) (RegisterIdentityResponse, error)
