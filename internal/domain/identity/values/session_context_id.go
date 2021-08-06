package identity

import merr "github.com/hywmongous/example-service/pkg/errors"

type SessionContextID string

func GenerateSessionContextID() SessionContextID {
	return SessionContextID(generateUuidValue())
}

func CreateSessionContextID(value string) (SessionContextID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return SessionContextID(""), merr.CreateFailedInvocation("CreateSessionContextID", err)
	}
	return SessionContextID(uuid), nil
}

func RecreateSessionContextID(value string) SessionContextID {
	return SessionContextID(value)
}

func (id SessionContextID) ToString() string {
	return uuidValue(id).toString()
}

func (id SessionContextID) Equals(other SessionContextID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
