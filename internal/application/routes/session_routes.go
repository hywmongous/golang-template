package routes

import (
	"github.com/hywmongous/example-service/internal/application/controllers"
	"github.com/hywmongous/example-service/internal/lib"
)

type SessionRoutes struct {
	handler    lib.RequestHandler
	controller controllers.SessionController
}

func SessionRoutesFactory(
	handler lib.RequestHandler,
	controller controllers.SessionController,
) SessionRoutes {
	return SessionRoutes{
		handler:    handler,
		controller: controller,
	}
}

func (routes SessionRoutes) Setup() {
	group := routes.handler.Gin.Group("/api/v1")
	// Get since we read all sessions for an account
	group.GET("/identities/:aid/sessions", routes.controller.GetAll)
	// POST since we create a new session for an account
	group.POST("/identities/:aid/sessions", routes.controller.Create)
	// PATCH since we partially updates all sessions by invalidating them
	group.PATCH("/identities/:aid/sessions", routes.controller.UseAll)
	// GET since we read a single session for an account
	group.GET("/identities/:aid/sessions/:sid", routes.controller.Get)
	// PATCH since we partially updates the session by invalidating it
	group.PATCH("/identities/:aid/sessions/:sid", routes.controller.Use)
}
