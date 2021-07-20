package application

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
	group := routes.handler.Gin.Group("/api/v1")
	// GET since we are reading all accounts
	group.GET("/accounts", routes.controller.GetAll)
	// POST since we are creating an account
	group.POST("/accounts", routes.controller.Create)
	// GET since we read a single account
	group.GET("/accounts/:aid", routes.controller.Get)
	// PATCH since we are patiallying updating an account
	group.PATCH("/accounts/:aid", routes.controller.Change)
	// DELETE since we are deleting the account (Aggregate root)
	group.DELETE("/accounts/:aid", routes.controller.Delete)
}
