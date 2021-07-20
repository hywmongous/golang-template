package bootstrap

import (
	application "github.com/hywmongous/example-service/internal/identity/application/controllers"

	"go.uber.org/fx"
)

var ControllerOptions = fx.Options(
	fx.Provide(application.AccountControllerFactory),
	fx.Provide(application.SessionControllerFactory),
	fx.Provide(application.TicketControllerFactory),
)
