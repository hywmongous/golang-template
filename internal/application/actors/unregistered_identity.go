package actors

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/application/usecases"
	"github.com/hywmongous/example-service/internal/domain/identity"
)

type UnregisteredIdentity struct {
	Register usecases.RegisterIdentityUseCase
}

func CreateUnregisteredIdentity() UnregisteredIdentity {
	return UnregisteredIdentity{
		Register,
	}
}

func Register(request usecases.RegisterIdentityRequest) (usecases.RegisterIdentityResponse, error) {
	identity, err := identity.CreateIdentity(
		request.Email,
		request.Password,
	)
	return usecases.RegisterIdentityResponse{
		IdentityID: identity.GetId(),
	}, errors.Wrap(err, "Something went wrong with the identity registration")
}
