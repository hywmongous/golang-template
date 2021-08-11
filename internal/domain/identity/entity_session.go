package identity

import (
	"errors"
)

type Session struct {
	id       SessionID
	revoked  bool
	contexts []SessionContext
}

var (
	ErrSessionRevoked = errors.New("session is revoked")
)

func CreateSession() Session {
	contexts := [1]SessionContext{CreateSessionContext()}
	return Session{
		id:       GenerateSessionID(),
		revoked:  false,
		contexts: contexts[:],
	}
}

func (session *Session) Refresh() SessionContext {
	CreateSessionContext := CreateSessionContext()
	session.contexts = append(session.contexts, CreateSessionContext)
	return CreateSessionContext
}

func (session *Session) Revoke() {
	session.revoked = true
}

func (session Session) Context() (SessionContext, error) {
	if session.revoked {
		return SessionContext{}, ErrSessionRevoked
	}

	latest := session.contexts[0]
	for _, context := range session.contexts {
		if context.issuedAt > latest.issuedAt {
			latest = context
		}
	}
	return latest, nil
}

func (session Session) GetId() SessionID {
	return session.id
}
