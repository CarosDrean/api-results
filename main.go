package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/helper"
	routes "github.com/CarosDrean/api-results.git/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	db.DB = helper.Get()
	r.HandleFunc("/", indexRouter)
	// r.HandleFunc("/api/login", middleware.Login)
	// r.HandleFunc("/validate", middleware.ValidateToken)
	s := r.PathPrefix("/api").Subrouter()
	routes.Routes(s)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:4800"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = helper.PORT //localhost
	}

	handler := c.Handler(r)

	fmt.Println("Server online!")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func indexRouter(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome!")
}
