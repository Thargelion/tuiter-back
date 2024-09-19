package security

import "net/http"

type AuthenticatorMiddleware struct {
	validator TokenValidator
}

func NewAuthenticator(validator TokenValidator) *AuthenticatorMiddleware {
	return &AuthenticatorMiddleware{validator: validator}
}

func (a *AuthenticatorMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		_, err := a.validator.ValidateToken(token)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
