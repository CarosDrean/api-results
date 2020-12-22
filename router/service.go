package routes

import (
	service "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func serviceRoutes(s *mux.Router) {
	s.HandleFunc("/all/{idProtocol}", mid.CheckSecurity(service.GetServicesPatientsWithProtocol)).Methods("GET")
	s.HandleFunc("/all-organization/{id}", mid.CheckSecurity(service.GetServicesPatientsWithOrganization)).Methods("GET")
	s.HandleFunc("/all-organization/", mid.CheckSecurity(service.GetServicesPatientsWithOrganizationFilter)).Methods("POST")
	s.HandleFunc("/filter/", mid.CheckSecurity(service.GetServicesPatientsWithProtocolFilter)).Methods("POST")
	s.HandleFunc("/filter-date/", mid.CheckSecurityInternalAdmin(service.GetServicesFilterDate)).Methods("POST")
}
