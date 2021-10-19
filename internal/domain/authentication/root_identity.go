package authentication

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/mediator"
)

var (
	ErrSessionNotFound   = errors.New("session could not be found")
	ErrIncorrectPassword = errors.New("password is incorrect")
)

type (
	IdentityID string
	Identity   struct {
		id       IdentityID
		email    Email
		password Password
		sessions []Session
	}
)

func (identity Identity) ID() IdentityID {
	return identity.id
}

func (identity *Identity) Email() Email {
	return identity.email
}

func (identity *Identity) Password() Password {
	return identity.password
}

func RecreateIdentity(
	id IdentityID,
	email Email,
	password Password,
	sessions []Session,
) Identity {
	return Identity{
		id:       id,
		email:    email,
		password: password,
		sessions: sessions,
	}
}

func Register(
	emailAddress string,
	plainTextPassword string,
) (Identity, error) {
	email, err := CreateEmail(emailAddress)
	if err != nil {
		return Identity{}, err
	}

	password, err := CreatePassword(plainTextPassword)
	if err != nil {
		return Identity{}, err
	}

	identity := Identity{
		id:       IdentityID(uuid.NewString()),
		email:    email,
		password: password,
		sessions: make([]Session, 0),
	}

	identity.publishEvent(&IdentityRegistered{
		ID:           string(identity.id),
		Email:        emailAddress,
		Passwordhash: password.hashedPassword,
	})

	return identity, nil
}

func (identity *Identity) Login(password string) (SessionID, error) {
	if err := identity.password.verify(password); err != nil {
		return SessionID(""), errors.Wrap(err, ErrIncorrectPassword.Error())
	}

	newSession, err := CreateSession()
	if err != nil {
		return SessionID(""), err
	}
	identity.sessions = append(identity.sessions, newSession)

	identity.publishEvent(&IdentityLoggedIn{
		SessionID: string(newSession.ID()),
	})

	return newSession.ID(), nil
}

func (identity *Identity) session(sessionID SessionID) (Session, error) {
	for _, session := range identity.sessions {
		if session.id == sessionID {
			return session, nil
		}
	}
	return Session{}, ErrSessionNotFound
}

func (identity *Identity) Logout(sessionID SessionID) error {
	session, err := identity.session(sessionID)
	if err != nil {
		return err
	}
	session.revoke()
	return nil
}

func (identity *Identity) publishEvent(event es.Data) {
	mediator.Publish(es.SubjectID(identity.Email().address), event)
}
