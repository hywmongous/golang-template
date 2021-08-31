package identity

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/pkg/crypto"
)

type ConfirmationHash string

var (
	ErrEmptyHash = errors.New("confirmation hash could not be create because the input value is empty")
)

func GenerateConfirmationHash() (ConfirmationHash, error) {
	str, err := crypto.GenerateRandomStringURLSafe(32)
	return ConfirmationHash(str), err
}

func CreateConfirmationHash(value string) (ConfirmationHash, error) {
	if value == "" {
		return "", ErrEmptyHash
	}
	return ConfirmationHash(value), nil
}

func recreateConfirmationHash(hash string) ConfirmationHash {
	return ConfirmationHash(hash)
}
