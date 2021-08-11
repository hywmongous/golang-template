package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type RegistrationRequest struct {
	Email    string
	Password string
}

type RegistrationResponse struct {
	IdentityID identity.IdentityID
}

type Register interface {
	DoRegister(request RegistrationRequest) (RegistrationResponse, error)
}
