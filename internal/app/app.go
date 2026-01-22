package app

import (
	"html/template"

	"github.com/fadlinrizqif/shortlink/internal/database"
)

type App struct {
	TemplateSignUp    *template.Template
	TemplateSignIn    *template.Template
	TemplateDashboard *template.Template
	DB                *database.Queries
	ServerSecret      string
}
