package routes

import (
	organization "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func organizationRoutes(s *mux.Router) {
	s.HandleFunc("/", mid.CheckSecurity(organization.GetOrganizations)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(organization.GetOrganization)).Methods("GET")
	s.HandleFunc("/send-mail", mid.RoleInternalAdmin(organization.SendURLTokenForExternalUser)).Methods("POST")
}
