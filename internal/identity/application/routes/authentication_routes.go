package identity

import (
	identity "github.com/hywmongous/example-service/internal/identity/application/controllers"
	"github.com/hywmongous/example-service/internal/lib"
)

type AuthenticationRoutes struct {
	handler    lib.RequestHandler
	controller identity.AuthenticationController
}

func AuthenticationRoutesFactory(
	handler lib.RequestHandler,
	controller identity.AuthenticationController,
) AuthenticationRoutes {
	return AuthenticationRoutes{
		handler:    handler,
		controller: controller,
	}
}

func (routes AuthenticationRoutes) Setup() {
	group := routes.handler.Gin.Group("/api/authentication")
	group.POST("/login", routes.controller.Login)
	group.POST("/logout", routes.controller.Login)
	group.POST("/refresh", routes.controller.Refresh)
}
