package values

import merr "github.com/hywmongous/example-service/pkg/errors"

type RefreshTokenID string

func GenerateRefreshTokenID() RefreshTokenID {
	return RefreshTokenID(generateUuidValue())
}

func CreateRefreshTokenID(value string) (RefreshTokenID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return RefreshTokenID(""), merr.CreateFailedInvocation("CreateRefreshTokenID", err)
	}
	return RefreshTokenID(uuid), nil
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
