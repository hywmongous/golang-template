package bootstrap

import (
	"context"

	"github.com/hywmongous/example-service/internal/lib"
	"go.uber.org/fx"
)

var Module = fx.Options(
	ControllerOptions,
	MiddlewareOptions,
	RouteOptions,
	LibOptions,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler lib.RequestHandler,
	routes Routes,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			go func() {
				routes.Setup()
				handler.Gin.Run(":8080")
			}()
			return nil
		},
		OnStop: func(c context.Context) error {
			return nil
		},
	})
}
