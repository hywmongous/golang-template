package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/presentation/gin/routes"
	"go.uber.org/fx"
)

func bootstrap(
	lifecycle fx.Lifecycle,
	engine *gin.Engine,
	routes routes.Routes,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			go func() {
				routes.Setup()
				if err := engine.Run(":5000"); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(c context.Context) error {
			return nil
		},
	})
}
