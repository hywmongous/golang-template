package identity

import (
	identity "github.com/hywmongous/example-service/internal/identity/application/controllers"
	"github.com/hywmongous/example-service/internal/lib"
)

type AccountRoutes struct {
	handler    lib.RequestHandler
	controller identity.AccountController
}

func AccountRoutesFactory(
	handler lib.RequestHandler,
	controller identity.AccountController,
) AccountRoutes {
	return AccountRoutes{
		handler:    handler,
		controller: controller,
	}
}

func (routes AccountRoutes) Setup() {
	group := routes.handler.Gin.Group("/api")
	group.GET("/accounts", routes.controller.GetAccount)
}
