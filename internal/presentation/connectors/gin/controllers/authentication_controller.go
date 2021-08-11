package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	identity_aggregate "github.com/hywmongous/example-service/internal/domain/identity"
	identity_values "github.com/hywmongous/example-service/internal/domain/identity"
	infrastructure "github.com/hywmongous/example-service/internal/infrastructure/services"
)

type AuthenticationController struct {
	jwtService infrastructure.JWTService
}

const (
	csrfHeaderKey             = "csrf"
	jwtAccessTokenCookieName  = "JWT-ACCESS-TOKEN"
	jwtRefreshTokenCookieName = "JWT-REFRESH-TOKEN"
)

var password, _ = identity_values.CreatePassword("password")
var email, _ = identity_values.CreateEmail("andreasbrandhoej@hotmail.com")
var currIdentity, _ = identity_aggregate.CreateIdentity(email, password)
var currSession, _ = currIdentity.Login("password")

func AuthenticationControllerFactory(
	jwtService infrastructure.JWTService,
) AuthenticationController {
	return AuthenticationController{
		jwtService: jwtService,
	}
}

func (controller AuthenticationController) Login(context *gin.Context) {
	currSession = identity_aggregate.CreateSession()
	username, password, ok := context.Request.BasicAuth()
	if username == "" || password == "" || !ok {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionContext, err := currSession.Context()
	if err != nil {
		context.Writer.WriteHeader(http.StatusInternalServerError)
		context.Abort()
	}

	controller.writeSessionToResponse(context, sessionContext)
	context.String(http.StatusOK, fmt.Sprintf("Logging in %s:%s", username, password))
}

func (controller AuthenticationController) Logout(context *gin.Context) {
	currSession.Revoke()
	context.String(http.StatusOK, "Logging out")
}

func (controller AuthenticationController) Refresh(context *gin.Context) {
	CreateSessionContext := currSession.Refresh()
	controller.writeSessionToResponse(context, CreateSessionContext)
	context.String(http.StatusOK, "Refresh session tokens")
}

func (controller AuthenticationController) Verify(context *gin.Context) {
	accessTokenCookie, err := context.Cookie(jwtAccessTokenCookieName)
	if err == nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshTokenCookie, err := context.Cookie(jwtRefreshTokenCookieName)
	if err == nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenPair := infrastructure.TokenPair{
		AccessToken:  accessTokenCookie,
		RefreshToken: refreshTokenCookie,
	}

	csrf := context.Request.Header.Get(csrfHeaderKey)

	if controller.jwtService.Verify(tokenPair, csrf, currSession) != nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	context.Writer.WriteHeader(http.StatusOK)
}

func (controller AuthenticationController) writeSessionToResponse(context *gin.Context, sessionContext identity_aggregate.SessionContext) {
	tokens, _ := controller.jwtService.Sign(currIdentity, sessionContext)

	context.Header(csrfHeaderKey, string(sessionContext.GetCsrf()))

	context.SetCookie(
		jwtAccessTokenCookieName,
		string(tokens.AccessToken),
		infrastructure.AccessTokenAbsoluteTimeoutDuration*60,
		"/",
		"localhost",
		false,
		true,
	)

	context.SetCookie(
		jwtRefreshTokenCookieName,
		string(tokens.RefreshToken),
		infrastructure.RefreshTokenAbsoluteTimeoutDuration*60,
		"/",
		"localhost",
		false,
		true,
	)
}
