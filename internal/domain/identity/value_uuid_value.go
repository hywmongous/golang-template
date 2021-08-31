package identity

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type uuidValue string

var (
	emptyUuidValue = uuidValue("")

	// ErrWrongFormat is returned when the creation of a uuid based on a value does not follow the format of uuids's
	ErrWrongFormat = errors.New("format is not a uuid")

	// ErrEmptyString is returned when the creation of uuid based on a value is the empty string
	ErrEmptyString = errors.New("value string is empty")
)

func generateUuidValue() uuidValue {
	return uuidValue(uuid.New().String())
}

func createUuidValue(value string) (uuidValue, error) {
	if value == "" {
		return emptyUuidValue, ErrEmptyString
	}
	parsedUuid, err := uuid.Parse(value)
	return uuidValue(string(parsedUuid[:])), errors.Wrap(err, "uuid parsing failed")
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
