package visitors

import (
	identity "github.com/hywmongous/example-service/internal/domain/identity"
	identity_values "github.com/hywmongous/example-service/internal/domain/identity"
)

type AggregateVisitor struct {
	root identity.Identity
}

func CreateAggregateVisitor() (AggregateVisitor, error) {
	return AggregateVisitor{
		root: identity.Identity{},
	}, nil
}

func RecreateAggregateVisitor(root identity.Identity) (AggregateVisitor, error) {
	return AggregateVisitor{
		root: root,
	}, nil
}

func (visitor AggregateVisitor) VisitRegisterIdentity(event identity.Registered) error {
	password := identity_values.RecreatePassword(event.GetHashedPassword())
	visitor.root = identity.RecreateIdentity(
		event.GetId(),
		identity_values.RecreateEmail(event.GetEmail(), false),
		password,
		[]identity.Session{},
		[]identity.Scope{},
	)
	return nil
}
