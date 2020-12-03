package routes

import (
	service "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func serviceRoutes(s *mux.Router) {
	s.HandleFunc("/all/{idProtocol}", mid.CheckSecurity(service.GetServicesPatientsWithProtocol)).Methods("GET")
}
