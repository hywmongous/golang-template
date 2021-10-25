package application

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/authentication"
	"github.com/hywmongous/example-service/internal/infrastructure"
	"github.com/hywmongous/example-service/internal/infrastructure/jaeger"
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

func (user UnregisteredUser) Register(
	ctx context.Context,
	request *RegisterIdentityRequest,
) (*RegisterIdentityResponse, error) {
	span, ctx := jaeger.StartSpanFromSpanContext(ctx, "Register")
	defer span.Finish()

	defer user.uow.Clear()

	identity, err := authentication.Register(
		request.Email,
		request.Password,
		user.uow.Mediator(),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrRegistrationFailed.Error())
	}

	if err = user.uow.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, ErrRegistrationFailedCommitting.Error())
	}

	return &RegisterIdentityResponse{
		Id: string(identity.ID()),
	}, nil
}
