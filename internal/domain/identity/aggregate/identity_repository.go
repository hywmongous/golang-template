package identity

import (
	values "github.com/hywmongous/example-service/internal/domain/identity/values"
)

type IdentityRepository interface {
	CreateIdentity() (Identity, error)
	ReadAllIdentities() ([]Identity, error)
	ReadIdentityById(id values.IdentityID) (Identity, error)
	// For CRQS it creates a new snapshot of the identity
	// The argument: Update occurs on the write mdoel (aggregate root)
	UpdateIdentity(identity Identity)
	DeleteIdentityById(id values.IdentityID) error
}
