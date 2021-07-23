package infrastructure

// Information:
// https://datatracker.ietf.org/doc/html/rfc7519#section-4
// https://curity.io/resources/learn/scopes-vs-claims/

import (
	"errors"
	"strings"
	"time"

	"github.com/hywmongous/example-service/internal/identity/domain"

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
	SessionId string   `json:"sid,omitempty"`
	Csrf      string   `json:"csrf,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	jwt.StandardClaims
}

var (
	ErrJwtSignatureInvalid = errors.New("JWT signature is invalid")
	ErrJwtHashUnavailable  = errors.New("JWT hash is unavailable")
	ErrJwtInvalidKeyType   = errors.New("JWT invalid key type for jwt service alg (HS256)")
	ErrJwtInvalidStructure = errors.New("JWT invalid token string structure")
	ErrJwtInvalidToken     = errors.New("JWT invalid token")

	ErrJwtVerificationSid  = errors.New("JWT session id is incorrect")
	ErrJwtVerificationCsrf = errors.New("JWT csrf is incorrect")
	ErrJwtVerificationJti  = errors.New("JWT id is incorrect")

	ErrJwtVerificationIssuedInTheFuture = errors.New("JWT IAT states it was issued in the future")
	ErrJwtVerificationTooEarly          = errors.New("JWT NBF states it cannot be used yet")
	ErrJwtVerificationExpired           = errors.New("JWT EXP states it has expired")

	ErrIncorrectCsrf = errors.New("CSRF is incorrect")
)

func JWTServiceFactory() JWTService {
	return JWTService{
		alg:        jwt.SigningMethodHS256,
		privateKey: []byte("Super secret string"),
	}
}

func createClaims(session domain.Session, token domain.Token) Claims {
	return Claims{
		SessionId: session.Id,
		Csrf:      string(session.Csrf[:]),
		StandardClaims: jwt.StandardClaims{
			Id:        token.Id,
			Issuer:    "hywmongous",
			IssuedAt:  token.IssuedAt,
			NotBefore: token.InitialTimeout,
			ExpiresAt: token.AbsoluteTimeout,
		},
	}
}

func (jwtService JWTService) Sign(session domain.Session) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(session, session.AccessToken),
	)
	accessTokenString, err := accessToken.SignedString(jwtService.privateKey)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(session, session.RefreshToken),
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

func (jwtService JWTService) Verify(token string) error {
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

func (jwtService JWTService) VerifyClaims(tokenPair TokenPair, csrf string, session domain.Session) error {
	if csrf != session.Csrf {
		return ErrIncorrectCsrf
	}

	accessTokenClaims, err := jwtService.Parse(tokenPair.AccessToken)
	if err != nil {
		return err
	}

	refreshTokenClaims, err := jwtService.Parse(tokenPair.RefreshToken)
	if err != nil {
		return err
	}

	if accessTokenClaims.SessionId != session.Id ||
		refreshTokenClaims.SessionId != session.Id {
		return ErrJwtVerificationCsrf
	}

	if accessTokenClaims.Csrf != session.Csrf ||
		refreshTokenClaims.Csrf != session.Csrf {
		return ErrJwtVerificationCsrf
	}

	if accessTokenClaims.Id != session.AccessToken.Id ||
		refreshTokenClaims.Id != session.RefreshToken.Id {
		return ErrJwtVerificationCsrf
	}

	now := time.Now().Unix()
	if accessTokenClaims.StandardClaims.IssuedAt > now ||
		refreshTokenClaims.StandardClaims.IssuedAt > now {
		return ErrJwtVerificationIssuedInTheFuture
	}

	if accessTokenClaims.StandardClaims.NotBefore > now ||
		refreshTokenClaims.StandardClaims.NotBefore > now {
		return ErrJwtVerificationTooEarly
	}

	if accessTokenClaims.StandardClaims.ExpiresAt > now ||
		refreshTokenClaims.StandardClaims.ExpiresAt > now {
		return ErrJwtVerificationExpired
	}

	return nil
}

func (jwtService JWTService) Parse(token string) (Claims, error) {
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
