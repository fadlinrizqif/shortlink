package middleware

import (
	"net/http"

	"github.com/fadlinrizqif/shortlink/internal/auth"
	"github.com/fadlinrizqif/shortlink/internal/database"
)

type AuthConfig struct {
	DB           *database.Queries
	ServerSecret string
}

func AuthMiddlware(next http.Handler, config *AuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getToken, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = auth.ValidateJWT(getToken.Value, config.ServerSecret)
		if err != nil {
			http.Redirect(w, r, "/auth/refresh", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
