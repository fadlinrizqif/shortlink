package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fadlinrizqif/shortlink/internal/app"
	"github.com/fadlinrizqif/shortlink/internal/database"
	"github.com/fadlinrizqif/shortlink/internal/handlers"
	"github.com/fadlinrizqif/shortlink/internal/middleware"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	serverSecret := os.Getenv("SECRET_SERVER")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	tmplSignUp := template.Must(template.ParseFiles("app/index.html"))
	tmplSignIn := template.Must(template.ParseFiles("app/login.html"))
	tmplDashboard := template.Must(template.ParseFiles("app/dashboard.html"))

	ApiConfig := &app.App{
		TemplateSignUp:    tmplSignUp,
		TemplateSignIn:    tmplSignIn,
		TemplateDashboard: tmplDashboard,
		DB:                dbQueries,
		ServerSecret:      serverSecret,
	}

	authConfig := middleware.AuthConfig{
		DB:           dbQueries,
		ServerSecret: serverSecret,
	}

	userHandler := handlers.NewUserHandlers(ApiConfig)
	linkHandler := handlers.NewLinkHandler(ApiConfig)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/signup", http.HandlerFunc(userHandler.CreateUsers))
	mux.HandleFunc("/login", http.HandlerFunc(userHandler.LoginUsers))
	mux.HandleFunc("/logout", http.HandlerFunc(userHandler.LogoutUsers))
	mux.Handle("/dashboard", middleware.AuthMiddlware(http.HandlerFunc(linkHandler.HandlerLink), &authConfig))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
