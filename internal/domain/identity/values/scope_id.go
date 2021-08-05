package values

import merr "github.com/hywmongous/example-service/pkg/errors"

type ScopeID string

func GenerateScopeID() ScopeID {
	return ScopeID(generateUuidValue())
}

func CreateScopeID(value string) (ScopeID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return ScopeID(""), merr.CreateFailedInvocation("CreateScopeID", err)
	}
	return ScopeID(uuid), nil
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
