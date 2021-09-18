package cqrs

import (
	"github.com/hywmongous/example-service/internal/domain/identity"
)

type identityModel struct {
	id       identity.IdentityID
	email    identity.Email
	password identity.Password
	sessions []identity.Session
}

func (model *identityModel) ApplyIdentityRegistered(event *identity.IdentityRegistered) readModel {
	model.id = identity.IdentityID(event.ID)
	model.email = identity.RecreateEmail(event.Email, false)
	model.password = identity.RecreatePassword(event.Passwordhash)
	return model
}

func (model *identityModel) ApplyIdentityLoggedIn(event *identity.IdentityLoggedIn) readModel {
	session := identity.RecreateSession(
		identity.SessionID(event.SessionID),
		false,
	)
	model.sessions = append(model.sessions, session)
	return model
}

func (model *identityModel) ApplyIdentityLoggedOut(event *identity.IdentityLoggedOut) readModel {
	sessionIndex := model.getSessionIndexById(identity.SessionID(event.SessionID))
	model.sessions[sessionIndex] = identity.RecreateSession(
		identity.SessionID(event.SessionID),
		true,
	)
	return model
}

func (model *identityModel) getSessionIndexById(id identity.SessionID) int {
	for idx, session := range model.sessions {
		if session.ID() == id {
			return idx
		}
	}
	return -1
}
