package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/application"
	"github.com/hywmongous/example-service/internal/infrastructure/jaeger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IdentityController struct {
	unregisteredUser application.UnregisteredUser
}

var ErrInvalidBasicAuth = errors.New("basic auth does not have email and password")

func AccountControllerFactory(
	unregisteredUser application.UnregisteredUser,
) IdentityController {
	return IdentityController{
		unregisteredUser: unregisteredUser,
	}
}

func (controller IdentityController) GetAll(context *gin.Context) {
	context.String(http.StatusOK, "Retrieving all accounts")
}

func (controller IdentityController) Create(context *gin.Context) {
	ctx := context.Request.Context()
	span := opentracing.SpanFromContext(ctx)

	email, password, ok := context.Request.BasicAuth()
	if email == "" || password == "" || !ok {
		context.String(http.StatusUnauthorized, "something went wrong with the basic auth")
		jaeger.SetError(span, ErrInvalidBasicAuth)

		return
	}

	span.LogFields(log.String("email", email))

	request := &application.RegisterIdentityRequest{
		Email:    email,
		Password: password,
	}

	response, err := controller.unregisteredUser.Register(ctx, request)
	if err != nil {
		context.String(http.StatusInternalServerError, err.Error())
		jaeger.SetError(span, err)

		return
	}

	context.JSON(http.StatusCreated, response)
}

func (controller IdentityController) Get(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Retrieving the account")
}

func (controller IdentityController) Change(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Changing the account information")
}

func (controller IdentityController) Delete(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Deleting the account")
}
