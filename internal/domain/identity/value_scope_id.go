package identity

import "github.com/cockroachdb/errors"

type ScopeID string

func GenerateScopeID() ScopeID {
	return ScopeID(generateUuidValue())
}

func CreateScopeID(value string) (ScopeID, error) {
	uuid, err := createUuidValue(value)
	return ScopeID(uuid), errors.Wrap(err, "createUuidValue")
}

func RecreateScopeID(value string) ScopeID {
	return ScopeID(value)
}

func (id ScopeID) ToString() string {
	return uuidValue(id).toString()
}

func (id ScopeID) Equals(other ScopeID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
