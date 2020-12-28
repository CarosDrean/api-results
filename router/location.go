package routes

import (
	location "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func locationRoutes(s *mux.Router) {
	ctrl := location.LocationController{DB: db.LocationDB{}}
	s.HandleFunc("/all/{idOrganization}", mid.CheckSecurity(ctrl.GetAllOrganizationID)).Methods("GET")
}
