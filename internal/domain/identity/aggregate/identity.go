package identity

import (
	"errors"

	values "github.com/hywmongous/example-service/internal/domain/identity/values"
)

type Identity struct {
	id       values.IdentityID
	email    values.Email
	password values.Password
	sessions []Session
	scopes   []Scope
}

var (
	ErrVerifyScopeNoHayMatches = errors.New("needle did not have a match in the haystack")

	ErrLogoutSessionNotFound       = errors.New("session id did not match any session")
	ErrLogoutSessionAlreadyRevoked = errors.New("attempted to logout from a revoked session")
)

func CreateIdentity(
	email values.Email,
	password values.Password,
) (Identity, error) {
	return Identity{
		id:       values.GenerateIdentityID(),
		email:    email,
		password: password,
	}, nil
}

func RecreateIdentity(
	id values.IdentityID,
	email values.Email,
	password values.Password,
	sessions []Session,
	scopes []Scope,
) Identity {
	return Identity{
		id:       id,
		email:    email,
		password: password,
		sessions: sessions,
		scopes:   scopes,
	}
}

func (identity *Identity) Login(password string) (Session, error) {
	if err := identity.password.Verify(password); err != nil {
		return Session{}, err
	}

	CreateSession := CreateSession()
	identity.sessions = append(identity.sessions, CreateSession)
	return CreateSession, nil
}

func (identity Identity) Logout(sessionId values.SessionID) error {
	var session Session
	var found bool
	for _, curr := range identity.sessions {
		if curr.GetId() == sessionId {
			session = curr
			found = true
			break
		}
	}

	if !found {
		return ErrLogoutSessionNotFound
	}

	if session.revoked {
		return ErrLogoutSessionAlreadyRevoked
	}

	session.Revoke()

	return nil
}

func (identity Identity) VerifyScope(scope string) error {
	needle, err := ParseScope(scope)
	if err != nil {
		return err
	}

	for _, hay := range identity.scopes {
		if HierarchicMatch(hay, needle) {
			return nil
		}
	}

	return ErrVerifyScopeNoHayMatches
}

func (identity Identity) GetId() values.IdentityID {
	return identity.id
}
