package values

import merr "github.com/hywmongous/example-service/pkg/errors"

type SessionID string

func GenerateSessionID() SessionID {
	return SessionID(generateUuidValue())
}

func CreateSessionID(value string) (SessionID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return SessionID(""), merr.CreateFailedInvocation("CreateSessionID", err)
	}
	return SessionID(uuid), nil
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
