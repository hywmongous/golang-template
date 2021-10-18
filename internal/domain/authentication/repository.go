package authentication

type Repository interface {
	FindIdentityByEmail(email string) (Identity, error)
}
