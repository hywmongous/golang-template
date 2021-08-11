package identity

type IdentityRepository interface {
	ReadAllIdentities() ([]Identity, error)
	ReadIdentityById(id IdentityID) (Identity, error)
	// For CRQS it creates a new snapshot of the identity
	// The argument: Update occurs on the write mdoel (aggregate root)
	UpdateIdentity(identity Identity) error
	DeleteIdentityById(id IdentityID) error
}
