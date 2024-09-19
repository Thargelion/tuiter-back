package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"tuiter.com/api/pkg/instant"
	"tuiter.com/api/pkg/syserror"
)

func NewJWTHandler(lifeSpan time.Duration, secret []byte, instant instant.Instant) *JWTHandler {
	return &JWTHandler{lifeSpan: lifeSpan, secret: secret, instant: instant}
}

type JWTHandler struct {
	lifeSpan time.Duration
	secret   []byte
	instant  instant.Instant
}

func (j *JWTHandler) GenerateToken(email string, username string) (string, error) {
	expirationTime := j.instant.Now().Add(j.lifeSpan)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "unlam-tuiter",
		"sub":   username,
		"email": email,
		"exp":   expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(j.secret)

	if err != nil {
		return "", fmt.Errorf("%w: %w", syserror.ErrInternal, err)
	}

	return tokenString, nil
}

func (j *JWTHandler) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %w", syserror.ErrUnauthorized, err)
	}

	return token, nil
}
