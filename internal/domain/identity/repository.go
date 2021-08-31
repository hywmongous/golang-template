package identity

type IdentityRepository interface {
	ReadIdentityById(id IdentityID) (Identity, error)
	RegisterIdentity(identity Identity) (Identity, error)
}
