package routes

import (
	"logger/config"
	"logger/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func GetMux(cnfg config.Config) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	handlers.New(cnfg.Models)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/log", handlers.WriteLog)

	return mux
}
