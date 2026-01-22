package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/fadlinrizqif/shortlink/internal/app"
	"github.com/fadlinrizqif/shortlink/internal/auth"
	"github.com/fadlinrizqif/shortlink/internal/database"
)

type UsersDetails struct {
	Name     string
	Email    string
	Password string
}

type UserHandler struct {
	App *app.App
}

type successVal struct {
	Success bool
	Message string
}

type errorVal struct {
	Error   bool
	Message string
}

func NewUserHandlers(app *app.App) *UserHandler {
	return &UserHandler{App: app}
}

func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.App.TemplateSignUp.Execute(w, nil)
		return
	}

	details := UsersDetails{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	hashPass, err := auth.HashPassword(details.Password)
	if err != nil {
		h.App.TemplateSignUp.Execute(w, errorVal{
			Error:   true,
			Message: "Something wrong in Server",
		})
		log.Fatal(err)
		return
	}

	userParam := database.CreateUserParams{
		Name:           details.Name,
		Email:          details.Email,
		HashedPassword: hashPass,
	}

	_, err = h.App.DB.CreateUser(r.Context(), userParam)
	if err != nil {
		h.App.TemplateSignUp.Execute(w, errorVal{
			Error:   true,
			Message: "Something wrong in Server",
		})
		log.Fatal(err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (h *UserHandler) LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.App.TemplateSignIn.Execute(w, nil)
		return
	}

	details := UsersDetails{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	getUser, err := h.App.DB.GetUser(r.Context(), details.Email)
	if err != nil {
		h.App.TemplateSignIn.Execute(w, errorVal{Error: true, Message: "User not found"})
		return
	}

	isPassCorrect, _ := auth.CheckPassword(details.Password, getUser.HashedPassword)
	if !isPassCorrect {
		h.App.TemplateSignIn.Execute(w, errorVal{Error: true, Message: "Wrong Password"})
		return
	}

	tokenDuration := time.Duration(60) * time.Minute
	jwtToken, err := auth.MakeJWT(getUser.ID, h.App.ServerSecret, tokenDuration)
	if err != nil {
		h.App.TemplateSignIn.Execute(w, errorVal{Error: true, Message: "Something wrong in server"})
		return
	}

	getRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		h.App.TemplateSignIn.Execute(w, errorVal{Error: true, Message: "Something wrong in server"})
		return
	}

	expireData := time.Now().AddDate(0, 0, 60)

	refreshParams := database.CreateRefreshTokenParams{
		Token:     getRefreshToken,
		UserID:    getUser.ID,
		ExpiresAt: expireData,
	}

	_, err = h.App.DB.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		h.App.TemplateSignIn.Execute(w, struct{ Error string }{"Something wrong in making jwt"})
		return
	}

	secure := r.TLS != nil

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    getRefreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 60,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) LogoutUsers(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.App.DB.UpdateRefreshToken(r.Context(), refreshToken.Value)
	}

	secure := r.TLS != nil

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
