package cqrs

import "github.com/hywmongous/example-service/internal/domain/identity"

type readModel interface {
	ApplyIdentityRegistered(event *identity.IdentityRegistered) readModel
	ApplyIdentityLoggedIn(event *identity.IdentityLoggedIn) readModel
	ApplyIdentityLoggedOut(event *identity.IdentityLoggedOut) readModel
}
