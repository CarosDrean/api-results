package routes

import (
	organization "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func organizationRoutes(s *mux.Router) {
	ctrl := organization.OrganizationController{DB: db.OrganizationDB{}}
	s.HandleFunc("/", mid.CheckSecurity(ctrl.GetAll)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Get)).Methods("GET")
	s.HandleFunc("/all-working-employer/{idUser}", mid.CheckSecurity(ctrl.GetAllWorkingOfEmployer)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Update)).Methods("PUT")
	s.HandleFunc("/send-mail", mid.RoleInternalAdmin(ctrl.SendURLTokenForExternalUser)).Methods("POST")
}
