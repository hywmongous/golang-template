package services

// Information:
// https://datatracker.ietf.org/doc/html/rfc7519#section-4
// https://curity.io/resources/learn/scopes-vs-claims/

import (
	"errors"
	"strings"
	"time"

	identity "github.com/hywmongous/example-service/internal/domain/identity/aggregate"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	alg        *jwt.SigningMethodHMAC
	privateKey []byte
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	SessionId string `json:"sid,omitempty"`
	Csrf      string `json:"csrf,omitempty"`
	jwt.StandardClaims
}

var (
	ErrJwtSignatureInvalid = errors.New("JWT signature is invalid")
	ErrJwtHashUnavailable  = errors.New("JWT hash is unavailable")
	ErrJwtInvalidKeyType   = errors.New("JWT invalid key type for jwt service alg (HS256)")
	ErrJwtInvalidStructure = errors.New("JWT invalid token string structure")
	ErrJwtInvalidToken     = errors.New("JWT invalid token")

	ErrVerificationNotIssuedAtTheSameTime = errors.New("the context was issued at a different time")
	ErrVerificationSessionRevoked         = errors.New("session is revoked")
	ErrVerificationIncorrestSessionId     = errors.New("session id is incorrect")
	ErrVerificationIncorrectCsrf          = errors.New("CSRF is incorrect")
	ErrVerificationIncorrectTokenId       = errors.New("the token id does not match the one in the context")
	ErrVerificationIncorrectIssueTime     = errors.New("the issue time is incorrect")
)

type Token struct {
	Id              string
	IssuedAt        int64
	InitialTimeout  int64
	AbsoluteTimeout int64
	Subject         string
}

const (
	AccessTokenInitialTimeoutDuration   = 0
	AccessTokenAbsoluteTimeoutDuration  = 30
	RefreshTokenInitialTimeoutDuration  = 15
	RefreshTokenAbsoluteTimeoutDuration = 30
)

func JWTServiceFactory() JWTService {
	return JWTService{
		alg:        jwt.SigningMethodHS256,
		privateKey: []byte("Super secret string"),
	}
}

func accessTokenFactory(identity identity.Identity, context identity.SessionContext) Token {
	// Access tokens can be used immediately and expires after 30 minutes
	now := time.Now()
	return Token{
		Subject:         identity.GetId().ToString(),
		Id:              context.GetAccessTokenId().ToString(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(AccessTokenInitialTimeoutDuration * time.Minute).Unix(),
		AbsoluteTimeout: now.Add(AccessTokenAbsoluteTimeoutDuration * time.Minute).Unix(),
	}
}

func refreshTokenFactory(identity identity.Identity, context identity.SessionContext) Token {
	// Refresh tokens can be used after 15 minutes and expires after 30
	now := time.Now()
	return Token{
		Subject:         identity.GetId().ToString(),
		Id:              context.GetRefreshTokenId().ToString(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(RefreshTokenInitialTimeoutDuration * time.Minute).Unix(),
		AbsoluteTimeout: now.Add(RefreshTokenAbsoluteTimeoutDuration * time.Minute).Unix(),
	}
}

func createClaims(context identity.SessionContext, token Token) Claims {
	return Claims{
		SessionId: context.GetId().ToString(),
		Csrf:      context.GetCsrf().ToString(),
		StandardClaims: jwt.StandardClaims{
			Id:        token.Id,
			Issuer:    "hywmongous",
			IssuedAt:  token.IssuedAt,
			NotBefore: token.InitialTimeout,
			ExpiresAt: token.AbsoluteTimeout,
		},
	}
}

func (jwtService JWTService) Sign(identity identity.Identity, context identity.SessionContext) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(context, accessTokenFactory(identity, context)),
	)
	accessTokenString, err := accessToken.SignedString(jwtService.privateKey)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(context, refreshTokenFactory(identity, context)),
	)
	refreshTokenString, err := refreshToken.SignedString(jwtService.privateKey)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (jwtService JWTService) Verify(tokenPair TokenPair, csrf string, session identity.Session) error {
	accessTokenClaims, err := jwtService.parse(tokenPair.AccessToken)
	if err != nil {
		session.Revoke()
		return err
	}

	refreshTokenClaims, err := jwtService.parse(tokenPair.RefreshToken)
	if err != nil {
		session.Revoke()
		return err
	}

	return jwtService.verifyClaims(accessTokenClaims, refreshTokenClaims, csrf, session)
}

func (JWTService JWTService) verifyClaims(accessToken Claims, refreshToken Claims, csrf string, session identity.Session) error {
	if accessToken.Csrf != csrf ||
		refreshToken.Csrf != csrf {
		return ErrVerificationIncorrectCsrf
	}

	if accessToken.SessionId != session.GetId().ToString() ||
		refreshToken.SessionId != session.GetId().ToString() {
		return ErrVerificationIncorrestSessionId
	}

	context, err := session.Context()
	if err != nil {
		return err
	}

	if accessToken.IssuedAt != context.GetIssuedAt().GetInt64() ||
		refreshToken.IssuedAt != context.GetIssuedAt().GetInt64() {
		return ErrVerificationIncorrectIssueTime
	}

	if accessToken.Id != context.GetAccessTokenId().ToString() ||
		refreshToken.Id != context.GetRefreshTokenId().ToString() {
		return ErrVerificationIncorrectTokenId
	}

	return nil
}

func (jwtService JWTService) parse(token string) (Claims, error) {
	if err := jwtService.verifyToken(token); err != nil {
		return Claims{}, err
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(
		token, claims,
		func(t *jwt.Token) (interface{}, error) { return jwtService.privateKey, nil },
	)

	switch err {
	case jwt.ErrSignatureInvalid:
	}

	switch err {
	case jwt.ErrSignatureInvalid:
		return *claims, ErrJwtSignatureInvalid
	case jwt.ErrHashUnavailable:
		return *claims, ErrJwtHashUnavailable
	}

	if !parsedToken.Valid {
		return *claims, ErrJwtInvalidToken
	}

	return *claims, nil
}

func (jwtService JWTService) verifyToken(token string) error {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return ErrJwtInvalidStructure
	}

	err := jwtService.alg.Verify(
		strings.Join(parts[0:2], "."),
		parts[2],
		jwtService.privateKey,
	)

	switch err {
	case jwt.ErrInvalidKey:
		return ErrJwtInvalidKeyType
	case jwt.ErrHashUnavailable:
		return ErrJwtHashUnavailable
	case jwt.ErrSignatureInvalid:
		return ErrJwtSignatureInvalid
	}

	return err
}
