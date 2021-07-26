package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hywmongous/example-service/internal/identity/domain"
	infrastructure "github.com/hywmongous/example-service/internal/identity/infrastructure/services"
)

type AuthenticationController struct {
	jwtService infrastructure.JWTService
}

const (
	csrfHeaderKey             = "csrf"
	jwtAccessTokenCookieName  = "JWT-ACCESS-TOKEN"
	jwtRefreshTokenCookieName = "JWT-REFRESH-TOKEN"
)

var session = domain.SessionFactory()

func AuthenticationControllerFactory(
	jwtService infrastructure.JWTService,
) AuthenticationController {
	return AuthenticationController{
		jwtService: jwtService,
	}
}

func (controller AuthenticationController) Login(context *gin.Context) {
	session = domain.SessionFactory()
	username, password, ok := context.Request.BasicAuth()
	if username == "" || password == "" || !ok {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionContext, err := session.Context()
	if err != nil {
		context.Writer.WriteHeader(http.StatusInternalServerError)
		context.Abort()
	}

	controller.writeSessionToResponse(context, sessionContext)
	context.String(http.StatusOK, fmt.Sprintf("Logging in %s:%s", username, password))
}

func (controller AuthenticationController) Logout(context *gin.Context) {
	session.Revoke()
	context.String(http.StatusOK, "Logging out")
}

func (controller AuthenticationController) Refresh(context *gin.Context) {
	newSessionContext := session.Refresh()
	controller.writeSessionToResponse(context, newSessionContext)
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

	if controller.jwtService.Verify(tokenPair, csrf, session) != nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	context.Writer.WriteHeader(http.StatusOK)
}

func (controller AuthenticationController) writeSessionToResponse(context *gin.Context, sessionContext domain.SessionContext) {
	tokens, _ := controller.jwtService.Sign(sessionContext)

	context.Header(csrfHeaderKey, string(sessionContext.Csrf))

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
