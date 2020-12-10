package routes

import (
	organization "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func organizationRoutes(s *mux.Router) {
	s.HandleFunc("/{id}", mid.CheckSecurity(organization.GetOrganization)).Methods("GET")
}
