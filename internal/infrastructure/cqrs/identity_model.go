package cqrs

import (
	"github.com/hywmongous/example-service/internal/domain/authentication"
)

type identityModel struct {
	id       authentication.IdentityID
	email    authentication.Email
	password authentication.Password
	sessions []authentication.Session
}

func defaultIdentityModel() identityModel {
	return identityModel{
		id:       authentication.IdentityID(""),
		email:    authentication.RecreateEmail("", false),
		password: authentication.RecreatePassword(""),
		sessions: make([]authentication.Session, 0),
	}
}

func (model *identityModel) ApplyIdentityRegistered(event *authentication.IdentityRegistered) readModel {
	model.id = authentication.IdentityID(event.ID)
	model.email = authentication.RecreateEmail(event.Email, false)
	model.password = authentication.RecreatePassword(event.Passwordhash)
	return model
}

func (model *identityModel) ApplyIdentityLoggedIn(event *authentication.IdentityLoggedIn) readModel {
	session := authentication.RecreateSession(
		authentication.SessionID(event.SessionID),
		false,
	)
	model.sessions = append(model.sessions, session)
	return model
}

func (model *identityModel) ApplyIdentityLoggedOut(event *authentication.IdentityLoggedOut) readModel {
	sessionIndex := model.getSessionIndexByID(authentication.SessionID(event.SessionID))
	model.sessions[sessionIndex] = authentication.RecreateSession(
		authentication.SessionID(event.SessionID),
		true,
	)
	return model
}

func (model *identityModel) getSessionIndexByID(id authentication.SessionID) int {
	for idx, session := range model.sessions {
		if session.ID() == id {
			return idx
		}
	}
	return -1
}
