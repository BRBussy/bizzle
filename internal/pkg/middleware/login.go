package middleware

import (
	"net/http"
)

type Login struct {
	preSharedSecret string
}

func (a *Login) Setup(
	preSharedSecret string,
) *Login {
	a.preSharedSecret = preSharedSecret
	return a
}

func (a *Login) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for pre-shared secret header
		pss := r.Header.Get("Pre-Shared-Secret")
		if pss == a.preSharedSecret {
			// if pre-shared secret present no authentication required
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
