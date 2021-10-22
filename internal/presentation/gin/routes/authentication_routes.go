package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/presentation/gin/controllers"
)

type AuthenticationRoutes struct {
	engine     *gin.Engine
	controller controllers.AuthenticationController
}

func CreateAuthenticationRoutes(
	engine *gin.Engine,
	controller controllers.AuthenticationController,
) AuthenticationRoutes {
	return AuthenticationRoutes{
		engine:     engine,
		controller: controller,
	}
}

func (routes AuthenticationRoutes) Setup() {
	group := routes.engine.Group("/api/v1/authentication")
	// PUT since we might come in a situation where we update
	// the session. Eg. if the authentication for a given session
	// has failed numerous times and we want to revoke the session.
	// We dont use POST because we are not always creating a new resource.
	// group.PUT("/verify")
	// POST since we are creating a session upon logging in an account
	group.POST("/login", routes.controller.Login)
	// POST since it is not idempotent and we update the session
	group.POST("/logout", routes.controller.Logout)
	// POST since we are creating a new tokens and csrf
	group.POST("/refresh", routes.controller.Refresh)
}
