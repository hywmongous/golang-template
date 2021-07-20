package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct{}

func AuthenticationControllerFactory() TicketController {
	return TicketController{}
}

func (controller AuthenticationController) Login(context *gin.Context) {
	username, password, ok := context.Request.BasicAuth()
	context.String(http.StatusOK, fmt.Sprintf("Logging ing %s:%s %t", username, password, ok))
}

func (controller AuthenticationController) Logout(context *gin.Context) {
	context.String(http.StatusOK, "Logging out")
}
