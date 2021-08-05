package values

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	merr "github.com/hywmongous/example-service/pkg/errors"
)

type uuidValue string

var (
	ErrUuidStringIsEmpty = errors.New("cannot be empty")
)

func generateUuidValue() uuidValue {
	return uuidValue(uuid.New().String())
}

func createUuidValue(value string) (uuidValue, error) {
	if value == "" {
		return "", merr.CreateInvalidInputError(
			"createUuidValue", "value", ErrUuidStringIsEmpty,
		)
	}
	// TODO: Verify format of uuid string value (Maybe by attempting a parse?)
	return uuidValue(value), nil
}

func recreateUuidValue(value string) uuidValue {
	return uuidValue(value)
}

func (id uuidValue) toString() string {
	return string(id)
}

func (id uuidValue) equals(other uuidValue) bool {
	return id.toString() == other.toString()
}
