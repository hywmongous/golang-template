package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/presentation/gin/controllers"
)

type TicketRoutes struct {
	controller controllers.TicketController
	engine     *gin.Engine
}

func CreateTicketRoutes(
	engine *gin.Engine,
	controller controllers.TicketController,
) TicketRoutes {
	return TicketRoutes{
		engine:     engine,
		controller: controller,
	}
}

func (routes TicketRoutes) Setup() {
	group := routes.engine.Group("/api/v1")
	// GET since we are reading tickets
	group.GET("/identities/:aid/tickets", routes.controller.GetAll)
	// POST since a new ticket is created
	group.POST("/identities/:aid/tickets", routes.controller.Create)
	// GET since we are reading a ticket
	group.GET("/identities/:aid/tickets/:tid", routes.controller.Get)
	// PATCH since we partially updates the ticket by invalidating it
	group.PATCH("/identities/:aid/tickets/:tid", routes.controller.Use)
}
