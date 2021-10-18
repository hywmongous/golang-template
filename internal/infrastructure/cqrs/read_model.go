package cqrs

import "github.com/hywmongous/example-service/internal/domain/authentication"

type readModel interface {
	ApplyIdentityRegistered(event *authentication.IdentityRegistered) readModel
	ApplyIdentityLoggedIn(event *authentication.IdentityLoggedIn) readModel
	ApplyIdentityLoggedOut(event *authentication.IdentityLoggedOut) readModel
}
