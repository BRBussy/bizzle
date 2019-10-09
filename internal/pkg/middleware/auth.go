package middleware

import (
	"fmt"
	"net/http"
)

type AuthMiddleware struct {
}

func (a *AuthMiddleware) Setup() *AuthMiddleware {
	return a
}

func (a *AuthMiddleware) ApplyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("middleware!!\n")
		next.ServeHTTP(w, r)
	})
}
