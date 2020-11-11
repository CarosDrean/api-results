package main

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/middleware"
	routes "github.com/CarosDrean/api-results.git/router"
	"github.com/gonutz/w32"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, w32.SW_HIDE)
		}
	}
	api()
}

func api(){
	r := mux.NewRouter()

	db.DB = helper.Get()

	r.HandleFunc("/", indexRouter)
	r.HandleFunc("/api/login", middleware.Login)
	r.HandleFunc("/file", controller.DownloadPDF)
	// r.HandleFunc("/validate", middleware.ValidateToken)
	s := r.PathPrefix("/api").Subrouter()
	routes.Routes(s)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://192.241.159.224"},
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
