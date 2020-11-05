package routes

import (
	file "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func fileRoutes(s *mux.Router) {
	s.HandleFunc("/", mid.CheckSecurity(file.DownloadPDF)).Methods("POST")
}
