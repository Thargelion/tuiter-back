package security

import "github.com/golang-jwt/jwt/v5"

type TokenHandler interface {
	GenerateToken(email string, username string) (string, error)
}

type TokenValidator interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}
