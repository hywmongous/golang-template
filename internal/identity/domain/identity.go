package domain

import (
	"github.com/hywmongous/example-service/pkg/guid"
)

type Identity struct {
	Id string
}

func IdentityFactory() Identity {
	return Identity{
		Id: guid.New().String(),
	}
}
