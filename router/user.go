package routes

import (
	user "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func userRoutes(s *mux.Router) {
	ctrl := user.UserController{}
	s.HandleFunc("/all/{id}", mid.CheckSecurity(ctrl.GetAllOrganization)).Methods("GET")
	s.HandleFunc("/", mid.CheckSecurity(user.GetSystemUsersPerson)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(user.GetSystemUser)).Methods("GET")
	s.HandleFunc("/password/{id}", mid.CheckSecurity(user.UpdatePasswordSystemUser)).Methods("PUT")
	s.HandleFunc("/", mid.RoleInternalAdminOrTempOrExternalAdmin(user.CreateSystemUser)).Methods("POST")
	s.HandleFunc("/{id}", mid.RoleInternalAdminOrTempOrExternalAdmin(user.UpdateSystemUser)).Methods("PUT")
	s.HandleFunc("/{id}", mid.RoleInternalAdminOrTempOrExternalAdmin(user.DeleteSystemUser)).Methods("DELETE")
}