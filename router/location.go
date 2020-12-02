package routes

import (
	location "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func locationRoutes(s *mux.Router) {
	s.HandleFunc("/all/{idOrganization}", mid.CheckSecurity(location.GetLocationsWidthOrganizationID)).Methods("GET")
}
