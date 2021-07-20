package application

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
	group := routes.handler.Gin.Group("/api/v1/authentication")
	// POST since we are creating a session upon logging in an account
	group.POST("/login", routes.controller.Login)
	// POST since it is not idempotent and we update the session
	group.POST("/logout", routes.controller.Logout)
}
