package identity

import "github.com/cockroachdb/errors"

type SessionContextID string

func GenerateSessionContextID() SessionContextID {
	return SessionContextID(generateUuidValue())
}

func CreateSessionContextID(value string) (SessionContextID, error) {
	uuid, err := createUuidValue(value)
	return SessionContextID(uuid), errors.Wrap(err, "createUuidValue")
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
