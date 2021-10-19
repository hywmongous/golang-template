package authentication

import "github.com/google/uuid"

type (
	SessionID string
	Session   struct {
		id      SessionID
		revoked bool
	}
)

func (session Session) ID() SessionID {
	return session.id
}

func (session Session) Revoked() bool {
	return session.revoked
}

func CreateSession() (Session, error) {
	return Session{
		id:      SessionID(uuid.NewString()),
		revoked: false,
	}, nil
}

func RecreateSession(
	id SessionID,
	revoked bool,
) Session {
	return Session{
		id:      id,
		revoked: revoked,
	}
}

func (session *Session) revoke() {
	session.revoked = true
}
