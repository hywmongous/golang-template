package usecases

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type LoginRequest struct {
	identityID values.IdentityID
	password   string
}

type LoginResponse struct {
	id string
}

type Login interface {
	DoLogin(request LoginRequest) (LoginResponse, error)
}
