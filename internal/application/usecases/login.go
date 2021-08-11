package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type LoginRequest struct {
	identityID identity.IdentityID
	password   string
}

type LoginResponse struct {
	id string
}

type Login interface {
	DoLogin(request LoginRequest) (LoginResponse, error)
}
