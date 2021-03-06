package services

// Information:
// https://datatracker.ietf.org/doc/html/rfc7519#section-4
// https://curity.io/resources/learn/scopes-vs-claims/

import (
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	SessionID string `json:"sid,omitempty"`
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
	ErrVerificationIncorrestSessionID     = errors.New("session id is incorrect")
	ErrVerificationIncorrectCsrf          = errors.New("CSRF is incorrect")
	ErrVerificationIncorrectTokenID       = errors.New("the token id does not match the one in the context")
	ErrVerificationIncorrectIssueTime     = errors.New("the issue time is incorrect")

	ErrSigningToken = errors.New("token could not be signed as jwt")
)

type Token struct {
	ID              string
	IssuedAt        int64
	InitialTimeout  int64
	AbsoluteTimeout int64
	Subject         string
}

const (
	AccessTokenInitialTimeoutMinutes   = 0
	AccessTokenAbsoluteTimeoutMinutes  = 30
	RefreshTokenInitialTimeoutMinutes  = 15
	RefreshTokenAbsoluteTimeoutMinutes = 30

	AccessTokenInitialTimeoutDuration   = AccessTokenInitialTimeoutMinutes * time.Minute
	AccessTokenAbsoluteTimeoutDuration  = AccessTokenAbsoluteTimeoutMinutes * time.Minute
	RefreshTokenInitialTimeoutDuration  = RefreshTokenInitialTimeoutMinutes * time.Minute
	RefreshTokenAbsoluteTimeoutDuration = RefreshTokenAbsoluteTimeoutMinutes * time.Minute

	Issuer = "hywmongous"
)

func JWTServiceFactory() JWTService {
	return JWTService{
		alg:        jwt.SigningMethodHS256,
		privateKey: []byte("Super secret string"),
	}
}

func createAccessToken(subject string) Token {
	// Access tokens can be used immediately and expires after 30 minutes
	now := time.Now()

	return Token{
		Subject:         subject,
		ID:              uuid.NewString(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(AccessTokenInitialTimeoutDuration).Unix(),
		AbsoluteTimeout: now.Add(AccessTokenAbsoluteTimeoutDuration).Unix(),
	}
}

func createRefreshToken(subject string) Token {
	// Refresh tokens can be used after 15 minutes and expires after 30
	now := time.Now()

	return Token{
		Subject:         subject,
		ID:              uuid.NewString(),
		IssuedAt:        now.Unix(),
		InitialTimeout:  now.Add(RefreshTokenInitialTimeoutDuration).Unix(),
		AbsoluteTimeout: now.Add(RefreshTokenAbsoluteTimeoutDuration).Unix(),
	}
}

func createClaims(token Token, sid string, csrf string) Claims {
	return Claims{
		SessionID: sid,
		Csrf:      csrf,
		StandardClaims: jwt.StandardClaims{
			Id:        token.ID,
			Subject:   token.Subject,
			Issuer:    Issuer,
			IssuedAt:  token.IssuedAt,
			NotBefore: token.InitialTimeout,
			ExpiresAt: token.AbsoluteTimeout,
			Audience:  "",
		},
	}
}

func (jwtService JWTService) Sign(
	subject string,
	sid string,
	csrf string,
) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(createAccessToken(subject), sid, csrf),
	)

	accessTokenString, err := accessToken.SignedString(jwtService.privateKey)
	if err != nil {
		return TokenPair{}, errors.Wrap(err, ErrSigningToken.Error())
	}

	refreshToken := jwt.NewWithClaims(
		jwtService.alg,
		createClaims(createRefreshToken(subject), sid, csrf),
	)

	refreshTokenString, err := refreshToken.SignedString(jwtService.privateKey)
	if err != nil {
		return TokenPair{}, errors.Wrap(err, ErrSigningToken.Error())
	}

	return TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (jwtService JWTService) Verify(token string, csrf string) (*Claims, error) {
	jsonPartsCount := 3

	// Verify general structure
	parts := strings.Split(token, ".")
	if len(parts) != jsonPartsCount {
		return nil, ErrJwtInvalidStructure
	}

	// Verify the jwt signature
	if err := jwtService.alg.Verify(
		strings.Join(parts[0:jsonPartsCount-1], "."),
		parts[jsonPartsCount-1],
		jwtService.privateKey,
	); err != nil {
		return nil, errors.Wrap(err, ErrJwtInvalidStructure.Error())
	}

	// Parse claims
	claims := Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return jwtService.privateKey, nil },
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrJwtInvalidToken.Error())
	}

	// Check the successfulness of the parsing
	if !parsedToken.Valid {
		return nil, errors.Wrap(parsedToken.Claims.Valid(), ErrJwtInvalidToken.Error())
	}

	// Verify claims
	if claims.Csrf != csrf {
		return nil, ErrVerificationIncorrectCsrf
	}

	return &claims, nil
}
