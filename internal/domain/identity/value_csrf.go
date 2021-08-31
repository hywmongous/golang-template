package identity

import "github.com/cockroachdb/errors"

type Csrf string

var (
	emptyCsrf = Csrf("")
)

func GenerateCsrf() Csrf {
	return Csrf(generateUuidValue())
}

func CreateCsrf(value string) (Csrf, error) {
	uuid, err := createUuidValue(value)
	return Csrf(uuid), errors.Wrap(err, "createUuidValue")
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
