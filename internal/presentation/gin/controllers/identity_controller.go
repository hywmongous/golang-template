package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/application/actors"
	"github.com/hywmongous/example-service/internal/application/usecases"
)

type IdentityController struct{}

func AccountControllerFactory() IdentityController {
	return IdentityController{}
}

func (controller IdentityController) GetAll(context *gin.Context) {
	context.String(http.StatusOK, "Retriving all accounts")
}

func (controller IdentityController) Create(context *gin.Context) {
	actor := actors.CreateUnregisteredIdentity()
	request := usecases.RegisterIdentityRequest{}
	response, err := actor.Register(request)
	if err != nil {
		context.String(http.StatusInternalServerError, "Something went wrong when registering the identity")
	}
	context.JSON(http.StatusCreated, response)
}

func (controller IdentityController) Get(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Retriving the account")
}

func (controller IdentityController) Change(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Changing the account information")
}

func (controller IdentityController) Delete(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Deleting the account")
}
