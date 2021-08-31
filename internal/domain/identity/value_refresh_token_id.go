package identity

import "github.com/cockroachdb/errors"

type RefreshTokenID string

var (
	emptyRefreshTokenID = RefreshTokenID("")
)

func GenerateRefreshTokenID() RefreshTokenID {
	return RefreshTokenID(generateUuidValue())
}

func CreateRefreshTokenID(value string) (RefreshTokenID, error) {
	uuid, err := createUuidValue(value)
	return RefreshTokenID(uuid), errors.Wrap(err, "createUuidValue")
}

func RecreateRefreshTokenID(value string) RefreshTokenID {
	return RefreshTokenID(value)
}

func (id RefreshTokenID) ToString() string {
	return uuidValue(id).toString()
}

func (id RefreshTokenID) Equals(other RefreshTokenID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
