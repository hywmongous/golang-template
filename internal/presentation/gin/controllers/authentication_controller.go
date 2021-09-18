package controllers

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hywmongous/example-service/internal/application"
	"github.com/hywmongous/example-service/internal/infrastructure/services"
)

type AuthenticationController struct {
	jwtService     services.JWTService
	registeredUser application.RegisteredUser
}

const (
	csrfHeaderKey             = "Csrf"
	jwtAccessTokenCookieName  = "JWT-ACCESS-TOKEN"
	jwtRefreshTokenCookieName = "JWT-REFRESH-TOKEN"
)

func AuthenticationControllerFactory(
	jwtService services.JWTService,
	registeredUser application.RegisteredUser,
) AuthenticationController {
	return AuthenticationController{
		jwtService:     jwtService,
		registeredUser: registeredUser,
	}
}

func (controller AuthenticationController) Login(context *gin.Context) {
	email, password, ok := context.Request.BasicAuth()
	if email == "" || password == "" || !ok {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := &application.LoginIdentityRequest{
		Email:    email,
		Password: password,
	}

	response, err := controller.registeredUser.Login(request)
	if err != nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = controller.writeSessionToResponse(context, email, response.SessionID); err != nil {
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	context.JSON(http.StatusOK, response)
}

func (controller AuthenticationController) Logout(context *gin.Context) {
	csrf := context.Request.Header.Get(csrfHeaderKey)

	accessToken, err := context.Cookie(jwtAccessTokenCookieName)
	if err != nil {
		context.String(http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := controller.jwtService.Verify(accessToken, csrf)
	if err != nil {
		context.String(http.StatusUnauthorized, errors.Wrap(err, csrf).Error())
		return
	}

	request := &application.LogoutIdentityRequest{
		Email:     claims.Subject,
		SessionID: claims.SessionId,
	}

	response, err := controller.registeredUser.Logout(request)
	if err != nil {
		context.String(http.StatusUnauthorized, err.Error())
		return
	}

	context.JSON(http.StatusOK, response)
}

func (controller AuthenticationController) Refresh(context *gin.Context) {
	context.String(http.StatusOK, "Refresh")
}

func (controller AuthenticationController) writeSessionToResponse(
	context *gin.Context,
	subject string,
	sid string,
) error {
	csrf := uuid.NewString()
	tokens, err := controller.jwtService.Sign(subject, sid, csrf)
	if err != nil {
		return err
	}

	context.Header(csrfHeaderKey, csrf)

	path := "/"
	domain := "localhost"
	secure := false
	httponly := true

	context.SetCookie(
		jwtAccessTokenCookieName,
		string(tokens.AccessToken),
		services.AccessTokenAbsoluteTimeoutMinutes*60,
		path,
		domain,
		secure,
		httponly,
	)

	context.SetCookie(
		jwtRefreshTokenCookieName,
		string(tokens.RefreshToken),
		services.RefreshTokenAbsoluteTimeoutMinutes*60,
		"/",
		"localhost",
		secure,
		httponly,
	)

	return nil
}
