package routes

import (
	user "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func userRoutes(s *mux.Router) {
	ctrl := user.UserController{DB: db.UserDB{}}
	s.HandleFunc("/all/{id}", mid.CheckSecurity(ctrl.GetAllOrganization)).Methods("GET")
	s.HandleFunc("/", mid.CheckSecurity(ctrl.GetAllPerson)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Get)).Methods("GET")
	s.HandleFunc("/password/{id}", mid.CheckSecurity(ctrl.UpdatePassword)).Methods("PUT")
	s.HandleFunc("/", mid.RoleInternalAdminOrTemp(ctrl.Create)).Methods("POST")
	s.HandleFunc("/{id}", mid.RoleInternalAdminOrTemp(ctrl.Update)).Methods("PUT")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Delete)).Methods("DELETE")
}