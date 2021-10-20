package application

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/authentication"
	"github.com/hywmongous/example-service/internal/infrastructure"
)

type RegisteredUser struct {
	uow infrastructure.UnitOfWork
}

var (
	ErrCouldNotFindIdentity   = errors.New("identity could not be found by email")
	ErrLoginFailed            = errors.New("identity login failed")
	ErrLoginFailedCommitting  = errors.New("identity login failed")
	ErrLogoutFailed           = errors.New("identity logout failed")
	ErrLogoutFailedCommitting = errors.New("identity logout failed")
)

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
		return nil, errors.Wrap(err, ErrCouldNotFindIdentity.Error())
	}

	sessionID, err := me.Login(request.Password)
	if err != nil {
		return nil, errors.Wrap(err, ErrLoginFailed.Error())
	}

	if err = user.uow.Commit(); err != nil {
		return nil, errors.Wrap(err, ErrLoginFailedCommitting.Error())
	}

	return &LoginIdentityResponse{
		SessionID: string(sessionID),
	}, nil
}

func (user RegisteredUser) Logout(request *LogoutIdentityRequest) (*LogoutIdentityResponse, error) {
	defer user.uow.Clear()

	me, err := user.uow.IdentityRepository().FindIdentityByEmail(request.Email)
	if err != nil {
		return nil, errors.Wrap(err, ErrCouldNotFindIdentity.Error())
	}

	err = me.Logout(authentication.SessionID(request.SessionID))
	if err != nil {
		return nil, errors.Wrap(err, ErrLoginFailed.Error())
	}

	if err = user.uow.Commit(); err != nil {
		return nil, errors.Wrap(err, ErrLogoutFailedCommitting.Error())
	}

	return &LogoutIdentityResponse{
		Revoked: true,
	}, nil
}
