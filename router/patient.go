package routes

import (
	patient "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func patientRoutes(s *mux.Router) {
	s.HandleFunc("/{id}", mid.CheckSecurity(patient.GetPatient)).Methods("GET")
}
