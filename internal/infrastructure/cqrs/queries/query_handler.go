package queries

import "github.com/hywmongous/example-service/internal/infrastructure/cqrs/commands"

type QueryHandler interface {
	commands.CommandHandler
}
