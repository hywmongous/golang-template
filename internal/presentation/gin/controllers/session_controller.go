package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionController struct{}

func SessionControllerFactory() SessionController {
	return SessionController{}
}

func (controller SessionController) GetAll(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Retrieving all sessions")
}

func (controller SessionController) Create(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Creating a session")
}

func (controller SessionController) UseAll(context *gin.Context) {
	// accountId := context.Param("aid")
	context.String(http.StatusOK, "Using all sessions")
}

func (controller SessionController) Get(context *gin.Context) {
	// accountId := context.Param("aid")
	// sessionId := context.Param("sid")
	context.String(http.StatusOK, "Retrieving the session")
}

func (controller SessionController) Use(context *gin.Context) {
	// accountId := context.Param("aid")
	// sessionId := context.Param("sid")
	context.String(http.StatusOK, "Using session")
}
