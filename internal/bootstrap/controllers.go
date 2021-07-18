package bootstrap

import (
	identity "github.com/hywmongous/example-service/internal/identity/application/controllers"

	"go.uber.org/fx"
)

var ControllerOptions = fx.Options(
	fx.Provide(identity.AccountControllerFactory),
	fx.Provide(identity.AuthenticationControllerFactory),
)
