package bootstrap

import (
	"github.com/hywmongous/example-service/internal/infrastructure/services"

	"go.uber.org/fx"
)

var InfrastructureServices = fx.Options(
	fx.Provide(services.JWTServiceFactory),
)
