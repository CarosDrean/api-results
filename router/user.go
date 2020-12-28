package routes

import (
	user "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func userRoutes(s *mux.Router) {
	s.HandleFunc("/", mid.CheckSecurity(user.GetSystemUsersPerson)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(user.GetSystemUser)).Methods("GET")
	s.HandleFunc("/password/{id}", mid.CheckSecurity(user.UpdatePasswordSystemUser)).Methods("PUT")
	s.HandleFunc("/", mid.RoleInternalAdminOrTemp(user.CreateSystemUser)).Methods("POST")
	s.HandleFunc("/{id}", mid.RoleInternalAdminOrTemp(user.UpdateSystemUser)).Methods("PUT")
	s.HandleFunc("/{id}", mid.RoleInternalAdminOrTemp(user.DeleteSystemUser)).Methods("DELETE")
}