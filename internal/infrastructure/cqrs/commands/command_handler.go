package commands

import "github.com/hywmongous/example-service/internal/domain/identity"

type CommandHandler interface {
	VisitRegistered(registration identity.Registered) error
}
