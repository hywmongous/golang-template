package identity

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct{}

func AccountControllerFactory() AccountController {
	return AccountController{}
}

func (controller AccountController) GetAccount(context *gin.Context) {
	context.String(http.StatusOK, "Get account!")
}
