package identity

type Repository interface {
	FindIdentityByEmail(email string) (Identity, error)
}
