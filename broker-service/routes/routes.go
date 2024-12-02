package routes

import (
	"broker/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rabbitmq/amqp091-go"
)

func GetMux(conn *amqp091.Connection) http.Handler {
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
	controller := handlers.NewController(conn)
	// create a route to check the health of server
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", controller.Broker)

	mux.Post("/handle", controller.HandleSubmission)
	return mux
}
