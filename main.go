package main

import (
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"os"

	"github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/middleware"
	routes "github.com/CarosDrean/api-results.git/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	configService()
	r := mux.NewRouter()

	db.DB = helper.Get()

	r.HandleFunc("/", indexRouter)
	r.HandleFunc("/api/login", middleware.Login)
	r.HandleFunc("/file", controller.DownloadPDF)
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

func configService() {
	svcConfig := &service.Config{
		Name:        "ApiResults",
		DisplayName: "Api Results",
		Description: "Este servicio permite la descarga de resultados de los examenes.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}

func indexRouter(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome!")
}
