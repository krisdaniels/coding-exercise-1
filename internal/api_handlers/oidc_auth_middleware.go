package api_handlers

import (
	"coding_exercise/internal/lib"
	"log"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	GetHandler(next http.Handler) http.Handler
}

type OidcAuthMiddleware struct {
	oidc_provider lib.OidcProvider
}

func NewOidcAuthMiddleware(oidc_provider lib.OidcProvider) AuthMiddleware {
	return &OidcAuthMiddleware{
		oidc_provider: oidc_provider,
	}
}

func (m *OidcAuthMiddleware) GetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get authorization header, 401
		auth_header := r.Header.Get("Authorization")
		if auth_header == "" {
			log.Println("authorization header missing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// check for bearer, 401
		if !strings.HasPrefix(auth_header, "Bearer") {
			log.Println("Bearer not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// extract token, 401
		auth_split := strings.Split(strings.TrimSpace(auth_header), " ")
		if len(auth_split) != 2 {
			log.Println("token missing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// validate token, 401
		token := auth_split[1]
		_, err := m.oidc_provider.ValidateToken(token)
		if err != nil {
			log.Printf("token validation error: %s\n", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// todo: read claims and put on context

		// if token ok call next
		next.ServeHTTP(w, r)
	})
}
