package security

import "github.com/golang-jwt/jwt/v5"

type TokenHandler interface {
	GenerateToken(id uint, email string, username string) (string, error)
}

type TokenValidator interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type TokenClaimsExtractor interface {
	ExtractClaims(token *jwt.Token) (jwt.MapClaims, error)
}

type UserExtractor interface {
	ExtractUserId(token *jwt.Token) (uint, error)
}
