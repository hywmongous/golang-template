package identity

import (
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedPassword string
}

const (
	cost = 14
)

func CreatePassword(password string) (Password, error) {
	createdPassword := Password{}
	if err := createdPassword.Set(password); err != nil {
		return Password{}, nil
	}
	return createdPassword, nil
}

func RecreatePassword(hashedPassword string) Password {
	return Password{
		hashedPassword: hashedPassword,
	}
}

func (password *Password) Set(input string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), cost)
	if err != nil {
		return err
	}
	password.hashedPassword = string(bytes)
	return nil
}

func (password Password) Verify(input string) error {
	return bcrypt.CompareHashAndPassword([]byte(password.hashedPassword), []byte(input))
}
