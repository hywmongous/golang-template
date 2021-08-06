package identity

import merr "github.com/hywmongous/example-service/pkg/errors"

type AccessTokenID string

func GenerateAccessTokenID() AccessTokenID {
	return AccessTokenID(generateUuidValue())
}

func CreateAccessTokenID(value string) (AccessTokenID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return AccessTokenID(""), merr.CreateFailedInvocation("CreateAccessTokenID", err)
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
