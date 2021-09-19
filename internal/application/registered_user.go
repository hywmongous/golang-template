package application

import (
	"github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/internal/infrastructure"
)

type RegisteredUser struct {
	uow infrastructure.UnitOfWork
}

func RegisteredUserFactory(
	uow infrastructure.UnitOfWork,
) RegisteredUser {
	return RegisteredUser{
		uow: uow,
	}
}

func (user RegisteredUser) Login(request *LoginIdentityRequest) (*LoginIdentityResponse, error) {
	defer user.uow.Clear()

	me, err := user.uow.IdentityRepository().FindIdentityByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	sessionID, err := me.Login(request.Password)
	if err != nil {
		return nil, err
	}

	user.uow.Commit()
	return &LoginIdentityResponse{
		SessionID: string(sessionID),
	}, nil
}

func (user RegisteredUser) Logout(request *LogoutIdentityRequest) (*LogoutIdentityResponse, error) {
	defer user.uow.Clear()

	me, err := user.uow.IdentityRepository().FindIdentityByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	err = me.Logout(identity.SessionID(request.SessionID))
	if err != nil {
		return nil, err
	}

	if err = user.uow.Commit(); err != nil {
		return nil, err
	}
	return &LogoutIdentityResponse{
		Revoked: true,
	}, nil
}
