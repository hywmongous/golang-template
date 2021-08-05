package actors

import "github.com/hywmongous/example-service/internal/application/usecases"

type Unregistered struct {
	register usecases.Register
}

func CreateUnregistered(
	register usecases.Register,
) (Unregistered, error) {
	return Unregistered{
		register: register,
	}, nil
}
