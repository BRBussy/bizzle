package middleware

import (
	"net/http"
)

type Authentication struct {
	preSharedSecret string
}

func (a *Authentication) Setup(
	preSharedSecret string,
) *Authentication {
	a.preSharedSecret = preSharedSecret
	return a
}

func (a *Authentication) Apply(next http.Handler) http.Handler {
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
