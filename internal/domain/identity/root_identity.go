package identity

import (
	"github.com/cockroachdb/errors"
)

type Identity struct {
	id       IdentityID
	email    Email
	password Password
	sessions []Session
	scopes   []Scope
}

var (
	ErrVerifyScopeNoHayMatches     = errors.New("needle did not have a match in the haystack")
	ErrLogoutSessionNotFound       = errors.New("session id did not match any session")
	ErrLogoutSessionAlreadyRevoked = errors.New("attempted to logout from a revoked session")
)

func CreateIdentity(
	email Email,
	password Password,
) (Identity, error) {
	return Identity{
		id:       GenerateIdentityID(),
		email:    email,
		password: password,
	}, nil
}

func RecreateIdentity(
	id IdentityID,
	email Email,
	password Password,
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

func (identity Identity) Logout(sessionId SessionID) error {
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

func (identity Identity) GetId() IdentityID {
	return identity.id
}
