package usecases

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type RegistrationRequest struct {
	Email    string
	Password string
}

type RegistrationResponse struct {
	IdentityID values.IdentityID
}

type Register interface {
	DoRegister(request RegistrationRequest) (RegistrationResponse, error)
}
