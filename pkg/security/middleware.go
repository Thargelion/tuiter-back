package security

import (
	"context"
	"net/http"
)

type TokenKey string

const (
	TokenMan TokenKey = "token"
)

type AuthenticatorMiddleware struct {
	validator TokenValidator
}

func NewAuthenticatorMiddleware(validator TokenValidator) *AuthenticatorMiddleware {
	return &AuthenticatorMiddleware{validator: validator}
}

func (a *AuthenticatorMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")

		token, err := a.validator.ValidateToken(tokenString)

		if err != nil {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(request.Context(), TokenMan, token)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
