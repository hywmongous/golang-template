package identity

import (
	"fmt"
)

type AccessTokenID string

func GenerateAccessTokenID() AccessTokenID {
	return AccessTokenID(generateUuidValue())
}

func CreateAccessTokenID(value string) (AccessTokenID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return AccessTokenID(""), fmt.Errorf("AccessTokenID creation failed: %w", err)
	}
	return AccessTokenID(uuid), nil
}

func RecreateAccessTokenID(value string) AccessTokenID {
	return AccessTokenID(value)
}

func (id AccessTokenID) ToString() string {
	return uuidValue(id).toString()
}

func (id AccessTokenID) Equals(other AccessTokenID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
