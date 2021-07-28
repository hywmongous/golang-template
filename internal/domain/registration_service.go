package domain

import (
	"github.com/hywmongous/example-service/internal/domain/identity"
)

type RegistrationForm struct {
	name     string
	password string
}

type RegistrationService interface {
	register(form RegistrationForm, repository identity.IdentityRepository) error
}
