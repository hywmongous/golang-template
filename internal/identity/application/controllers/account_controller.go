package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct{}

func AccountControllerFactory() AccountController {
	return AccountController{}
}

func (controller AccountController) GetAll(context *gin.Context) {
	context.String(http.StatusOK, "Retriving all accounts")
}

func (controller AccountController) Create(context *gin.Context) {
	context.String(http.StatusOK, "Creating an account")
}

func (controller AccountController) Get(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Retriving the account")
}

func (controller AccountController) Change(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Changing the account information")
}

func (controller AccountController) Delete(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Deleting the account")
}
