package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/presentation/gin/controllers"
)

type AccountRoutes struct {
	engine     *gin.Engine
	controller controllers.IdentityController
}

func CreateAccountRoutes(
	engine *gin.Engine,
	controller controllers.IdentityController,
) AccountRoutes {
	return AccountRoutes{
		engine:     engine,
		controller: controller,
	}
}

func (routes AccountRoutes) Setup() {
	group := routes.engine.Group("/api/v1")
	// GET since we are reading all accounts
	group.GET("/identities", routes.controller.GetAll)
	// POST since we are creating an account
	group.POST("/identities", routes.controller.Create)
	// GET since we read a single account
	group.GET("/identities/:aid", routes.controller.Get)
	// PATCH since we are patiallying updating an account
	group.PATCH("/identities/:aid", routes.controller.Change)
	// DELETE since we are deleting the account (Aggregate root)
	group.DELETE("/identities/:aid", routes.controller.Delete)
}
