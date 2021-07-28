package identity

import (
	"errors"

	"github.com/hywmongous/example-service/pkg/guid"
)

type Session struct {
	Id       string
	revoked  bool
	Contexts []SessionContext
}

var (
	ErrSessionRevoked = errors.New("session is revoked")
)

func SessionFactory() Session {
	contexts := [1]SessionContext{SessionContextFactory()}
	return Session{
		Id:       guid.New().String(),
		revoked:  false,
		Contexts: contexts[:],
	}
}

func (session *Session) Refresh() SessionContext {
	newSessionContext := SessionContextFactory()
	session.Contexts = append(session.Contexts, newSessionContext)
	return newSessionContext
}

func (session *Session) Revoke() {
	session.revoked = true
}

func (session Session) Context() (SessionContext, error) {
	if session.revoked {
		return SessionContext{}, ErrSessionRevoked
	}

	latest := session.Contexts[0]
	for _, context := range session.Contexts {
		if context.IssuedAt > latest.IssuedAt {
			latest = context
		}
	}
	return latest, nil
}
