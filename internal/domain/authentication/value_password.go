package authentication

import (
	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedPassword string
}

const (
	cost = 14
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidPassword   = errors.New("password is invalid and could not be set")
)

func (password Password) HashedPassword() string {
	return password.hashedPassword
}

func CreatePassword(password string) (Password, error) {
	createdPassword := DefaultPassword()
	if err := createdPassword.set(password); err != nil {
		return DefaultPassword(), errors.Wrap(err, ErrInvalidPassword.Error())
	}

	return createdPassword, nil
}

func RecreatePassword(hashedPassword string) Password {
	return Password{
		hashedPassword: hashedPassword,
	}
}

func DefaultPassword() Password {
	return Password{
		hashedPassword: "",
	}
}

func (password *Password) set(input string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), cost)
	if err != nil {
		return errors.Wrap(
			err,
			ErrInvalidPassword.Error(),
		)
	}

	password.hashedPassword = string(bytes)

	return nil
}

func (password Password) verify(input string) error {
	return errors.Wrap(
		bcrypt.CompareHashAndPassword([]byte(password.hashedPassword), []byte(input)),
		ErrIncorrectPassword.Error(),
	)
}
