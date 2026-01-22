package handlers

import (
	"net/http"

	"github.com/fadlinrizqif/shortlink/internal/app"
)

type LinkHandler struct {
	LinkTmpl *app.App
}

func NewLinkHandler(app *app.App) *LinkHandler {
	return &LinkHandler{LinkTmpl: app}
}

func (h *LinkHandler) HandlerLink(w http.ResponseWriter, r *http.Request) {
	h.LinkTmpl.TemplateDashboard.Execute(w, nil)
}
