package values

import merr "github.com/hywmongous/example-service/pkg/errors"

type Csrf string

func GenerateCsrf() Csrf {
	return Csrf(generateUuidValue())
}

func CreateCsrf(value string) (Csrf, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return Csrf(""), merr.CreateFailedInvocation("CreateCsrf", err)
	}
	return Csrf(uuid), nil
}

func RecreateCsrf(value string) Csrf {
	return Csrf(value)
}

func (id Csrf) ToString() string {
	return uuidValue(id).toString()
}

func (id Csrf) Equals(other Csrf) bool {
	return uuidValue(id).equals(uuidValue(other))
}
