package application

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/authentication"
	"github.com/hywmongous/example-service/internal/infrastructure"
)

type UnregisteredUser struct {
	uow infrastructure.UnitOfWork
}

var (
	ErrRegistrationFailed           = errors.New("identity registration failed")
	ErrRegistrationFailedCommitting = errors.New("identity registration failed committing")
)

func UnregisteredUserFactory(
	uow infrastructure.UnitOfWork,
) UnregisteredUser {
	return UnregisteredUser{
		uow: uow,
	}
}

func (user UnregisteredUser) Register(request *RegisterIdentityRequest) (*RegisterIdentityResponse, error) {
	defer user.uow.Clear()

	identity, err := authentication.Register(
		request.Email,
		request.Password,
		user.uow.Mediator(),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrRegistrationFailed.Error())
	}

	if err = user.uow.Commit(); err != nil {
		return nil, errors.Wrap(err, ErrRegistrationFailedCommitting.Error())
	}

	return &RegisterIdentityResponse{
		Id: string(identity.ID()),
	}, nil
}
