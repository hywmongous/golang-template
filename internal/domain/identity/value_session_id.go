package identity

import "github.com/cockroachdb/errors"

type SessionID string

func GenerateSessionID() SessionID {
	return SessionID(generateUuidValue())
}

func CreateSessionID(value string) (SessionID, error) {
	uuid, err := createUuidValue(value)
	return SessionID(uuid), errors.Wrap(err, "createUuidValue")
}

func RecreateSessionID(value string) SessionID {
	return SessionID(value)
}

func (id SessionID) ToString() string {
	return uuidValue(id).toString()
}

func (id SessionID) Equals(other SessionID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
