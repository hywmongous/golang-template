package identity

import (
	"errors"

	"github.com/hywmongous/example-service/pkg/guid"
)

type Identity struct {
	Id       string
	password Password
	sessions []Session
	scopes   []Scope
}

var (
	ErrVerifyScopeNoHayMatches = errors.New("needle did not have a match in the haystack")

	ErrLogoutSessionNotFound       = errors.New("session id did not match any session")
	ErrLogoutSessionAlreadyRevoked = errors.New("attempted to logout from a revoked session")
)

func IdentityFactory() Identity {
	return Identity{
		Id: guid.New().String(),
	}
}

func (identity *Identity) Login(password string) (Session, error) {
	if err := identity.password.Verify(password); err != nil {
		return Session{}, err
	}

	newSession := SessionFactory()
	identity.sessions = append(identity.sessions, newSession)
	return newSession, nil
}

func (identity Identity) Logout(sessionId string) error {
	var session Session
	var found bool
	for _, curr := range identity.sessions {
		if curr.Id == sessionId {
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
