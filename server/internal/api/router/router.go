package router

import (
	"net/http"

	"github.com/example/go-webdb-template/internal/api/handler"
	"github.com/example/go-webdb-template/internal/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// NewRouter は新しいルーターを作成
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler, cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Post routes
	api.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts", postHandler.ListPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", postHandler.GetPost).Methods("GET")
	api.HandleFunc("/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	api.HandleFunc("/posts/{id}", postHandler.DeletePost).Methods("DELETE")

	// User-Post JOIN route
	api.HandleFunc("/user-posts", postHandler.GetUserPosts).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	})

	return c.Handler(r)
}
