package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketController struct{}

func TicketControllerFactory() TicketController {
	return TicketController{}
}

func (controller TicketController) Get(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Retrieving the ticket")
}

func (controller TicketController) Create(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Creating a ticket")
}

func (controller TicketController) GetAll(context *gin.Context) {
	// accountId := context.Param("aid")
	// ticketId := context.Params("tid")
	context.String(http.StatusOK, "Retrieving all tickets")
}

func (controller TicketController) Use(context *gin.Context) {
	// accountId := context.Param("aid")
	// ticketId := context.Params("tid")
	context.String(http.StatusOK, "Using the ticket")
}
