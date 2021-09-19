package application

import (
	"github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/internal/infrastructure"
)

type UnregisteredUser struct {
	uow infrastructure.UnitOfWork
}

func UnregisteredUserFactory(
	uow infrastructure.UnitOfWork,
) UnregisteredUser {
	return UnregisteredUser{
		uow: uow,
	}
}

func (user UnregisteredUser) Register(request *RegisterIdentityRequest) (*RegisterIdentityResponse, error) {
	defer user.uow.Clear()

	identity, err := identity.Register(
		request.Email,
		request.Password,
	)
	if err != nil {
		return nil, err
	}

	if err = user.uow.Commit(); err != nil {
		return nil, err
	}
	return &RegisterIdentityResponse{
		Id: string(identity.ID()),
	}, nil
}
