package routes

import (
	"mailer-service/email"
	"mailer-service/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func GetMux(m email.Mail) http.Handler {
	mux := chi.NewRouter()

	// define the specifications about the server
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://*",
			"https://*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           30,
	}))

	// create a route to check the health of server
	mux.Use(middleware.Heartbeat("/ping"))
	// init mail variable in package handlers
	handlers.New(m)
	mux.Post("/send", handlers.SendMail)
	return mux
}
