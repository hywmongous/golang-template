package identity

type IdentityRepository interface {
	insertIdentity(identity Identity) error
	getById(id string) (Identity, error)
}
