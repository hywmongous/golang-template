package identity

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct{}

func AuthenticationControllerFactory() AuthenticationController {
	return AuthenticationController{}
}

func (controller AuthenticationController) Login(context *gin.Context) {
	username, password, ok := context.Request.BasicAuth()
	context.String(http.StatusOK, fmt.Sprintf("Authorization %s:%s %t", username, password, ok))
}

func (controller AuthenticationController) Logout(context *gin.Context) {
	context.String(http.StatusOK, "Logout")
}

func (controller AuthenticationController) Refresh(context *gin.Context) {
	context.String(http.StatusOK, "Refresh")
}
