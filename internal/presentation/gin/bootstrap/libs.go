package bootstrap

import (
	"github.com/hywmongous/example-service/internal/lib"

	"go.uber.org/fx"
)

var LibOptions = fx.Options(
	fx.Provide(lib.RequestHandlerFactory),
)
