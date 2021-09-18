package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/application"
)

type IdentityController struct {
	unregisteredUser application.UnregisteredUser
}

func AccountControllerFactory(
	unregisteredUser application.UnregisteredUser,
) IdentityController {
	return IdentityController{
		unregisteredUser: unregisteredUser,
	}
}

func (controller IdentityController) GetAll(context *gin.Context) {
	context.String(http.StatusOK, "Retriving all accounts")
}

func (controller IdentityController) Create(context *gin.Context) {
	email, password, ok := context.Request.BasicAuth()
	if email == "" || password == "" || !ok {
		context.String(http.StatusUnauthorized, "something went wrong with the basic auth")
		return
	}

	request := &application.RegisterIdentityRequest{
		Email:    email,
		Password: password,
	}

	response, err := controller.unregisteredUser.Register(request)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		return
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
