package bootstrap

import (
	infrastructure "github.com/hywmongous/example-service/internal/identity/infrastructure/services"

	"go.uber.org/fx"
)

var InfrastructureServices = fx.Options(
	fx.Provide(infrastructure.JWTServiceFactory),
)
