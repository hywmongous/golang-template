package application

import (
	identity "github.com/hywmongous/example-service/internal/identity/application/controllers"
	"github.com/hywmongous/example-service/internal/lib"
)

type TicketRoutes struct {
	handler    lib.RequestHandler
	controller identity.TicketController
}

func TicketRoutesFactory(
	handler lib.RequestHandler,
	controller identity.TicketController,
) TicketRoutes {
	return TicketRoutes{
		handler:    handler,
		controller: controller,
	}
}

func (routes TicketRoutes) Setup() {
	group := routes.handler.Gin.Group("/api/v1")
	// GET since we are reading tickets
	group.GET("/accounts/:aid/tickets", routes.controller.GetAll)
	// POST since a new ticket is created
	group.POST("/accounts/:aid/tickets", routes.controller.Create)
	// GET since we are reading a ticket
	group.GET("/accounts/:aid/tickets/:tid", routes.controller.Get)
	// PATCH since we partially updates the ticket by invalidating it
	group.PATCH("/accounts/:aid/tickets/:tid", routes.controller.Use)
}
