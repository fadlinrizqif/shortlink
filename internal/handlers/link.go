package handlers

import (
	"net/http"

	"github.com/fadlinrizqif/shortlink/internal/app"
	"github.com/fadlinrizqif/shortlink/internal/auth"
	"github.com/fadlinrizqif/shortlink/internal/database"
)

type LinkHandler struct {
	App *app.App
}

type linkDetail struct {
	OriginalURL string
}

type returnVal struct {
	ShortURL bool
	NewLink  string
}

func NewLinkHandler(app *app.App) *LinkHandler {
	return &LinkHandler{App: app}
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + r.Host + "/"
}

type DashboardData struct {
	Error    bool
	ShortURL bool
	NewLink  string
	Message  string
}

func (h *LinkHandler) HandlerLink(w http.ResponseWriter, r *http.Request) {
	data := DashboardData{}

	if r.Method == http.MethodPost {
		details := linkDetail{
			OriginalURL: r.FormValue("original_url"),
		}
		newLinkCode, _ := auth.MakeLinkCode()

		getJWT, _ := r.Cookie("access_token")

		getUserID, err := auth.ValidateJWT(getJWT.Value, h.App.ServerSecret)
		if err != nil {
			h.App.TemplateSignIn.Execute(w, errorVal{Error: true, Message: "Something wrong in server"})
			return
		}

		_, err = h.App.DB.CreateNewLink(r.Context(), database.CreateNewLinkParams{
			Code:    newLinkCode,
			LinkUrl: details.OriginalURL,
			UserID:  getUserID,
		})

		if err != nil {
			h.App.TemplateDashboard.Execute(w, nil)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	getCookieToken, _ := r.Cookie("access_token")

	getUserID, _ := auth.ValidateJWT(getCookieToken.Value, h.App.ServerSecret)

	getNewCode, err := h.App.DB.GetNewLink(r.Context(), getUserID)
	if err != nil {
		h.App.TemplateDashboard.Execute(w, nil)
		return
	}

	data.ShortURL = true
	data.NewLink = baseURL(r) + getNewCode
	h.App.TemplateDashboard.Execute(w, data)

}

func (h *LinkHandler) RedirectLink(w http.ResponseWriter, r *http.Request) {
	codeLink := r.PathValue("codeLink")

	getLink, err := h.App.DB.GetLink(r.Context(), codeLink)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	http.Redirect(w, r, getLink.LinkUrl, http.StatusSeeOther)

}
