package security

import (
	"context"
	"net/http"
)

type AuthenticatorMiddleware struct {
	validator TokenValidator
}

func NewAuthenticatorMiddleware(validator TokenValidator) *AuthenticatorMiddleware {
	return &AuthenticatorMiddleware{validator: validator}
}

func (a *AuthenticatorMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		token, err := a.validator.ValidateToken(tokenString)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
