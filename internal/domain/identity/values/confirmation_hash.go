package identity

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/pkg/crypto"
	merr "github.com/hywmongous/example-service/pkg/errors"
)

type ConfirmationHash string

var (
	ErrHashEmpty = errors.New("hash string value is empty")
)

func GenerateConfirmationHash() (ConfirmationHash, error) {
	str, err := crypto.GenerateRandomStringURLSafe(32)
	return ConfirmationHash(str), err
}

func CreateConfirmationHash(value string) (ConfirmationHash, error) {
	if value == "" {
		return "", merr.CreateInvalidInputError(
			"CreateConfirmationHash", "value", ErrHashEmpty,
		)
	}
	return ConfirmationHash(value), nil
}

func recreateConfirmationHash(hash string) ConfirmationHash {
	return ConfirmationHash(hash)
}
