package bootstrap

import (
	"github.com/hywmongous/example-service/internal/application/controllers"

	"go.uber.org/fx"
)

var ControllerOptions = fx.Options(
	fx.Provide(controllers.AccountControllerFactory),
	fx.Provide(controllers.AuthenticationControllerFactory),
	fx.Provide(controllers.SessionControllerFactory),
	fx.Provide(controllers.TicketControllerFactory),
)
