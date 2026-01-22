package handlers

import (
	"net/http"
	"time"

	"github.com/fadlinrizqif/shortlink/internal/app"
	"github.com/fadlinrizqif/shortlink/internal/auth"
)

type AuthHandler struct {
	App *app.App
}

func NewAuthHandler(app *app.App) *AuthHandler {
	return &AuthHandler{App: app}
}

func (h *AuthHandler) CheckRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, _ := r.Cookie("refresh_token")

	refreshTokenDB, err := h.App.DB.GetRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenDuration := time.Duration(60) * time.Minute

	newToken, err := auth.MakeJWT(refreshTokenDB.UserID, h.App.ServerSecret, tokenDuration)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	secure := r.TLS != nil

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
