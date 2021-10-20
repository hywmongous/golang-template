package controllers

import (
	"log"
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
	csrfHeaderKey = "Csrf"
	/* #nosec */
	jwtAccessTokenCookieName = "JWT-ACCESS-TOKEN"
	/* #nosec */
	jwtRefreshTokenCookieName = "JWT-REFRESH-TOKEN"
	secondsPerMinute          = 60
)

var ErrCouldNotWriteSessionToResponse = errors.New("session could not be written to response")

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
		log.Println("Login endpoint error", err)
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = controller.writeSessionToResponse(context, email, response.SessionID); err != nil {
		log.Println("Login endpoint error", err)
		context.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	if response == nil {
		log.Println("Login response was nil")
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
		SessionID: claims.SessionID,
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
		return errors.Wrap(
			err,
			ErrCouldNotWriteSessionToResponse.Error(),
		)
	}

	context.Header(csrfHeaderKey, csrf)

	path := "/"
	domain := "localhost"
	secure := false
	httponly := true

	context.SetCookie(
		jwtAccessTokenCookieName,
		tokens.AccessToken,
		services.AccessTokenAbsoluteTimeoutMinutes*secondsPerMinute,
		path,
		domain,
		secure,
		httponly,
	)

	context.SetCookie(
		jwtRefreshTokenCookieName,
		tokens.RefreshToken,
		services.RefreshTokenAbsoluteTimeoutMinutes*secondsPerMinute,
		"/",
		"localhost",
		secure,
		httponly,
	)

	return nil
}
