package routes

import (
	user "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func userRoutes(s *mux.Router) {
	// s.HandleFunc("/", user.GetUsers).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(user.GetSystemUser)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(user.UpdatePasswordSystemUser)).Methods("PUT")
	// s.HandleFunc("/", user.CreateUser).Methods("POST")
	// s.HandleFunc("/{id}", mid.CheckSecurity(user.UpdateUser)).Methods("PUT")
	// s.HandleFunc("/{id}", mid.CheckSecurity(user.DeleteUser)).Methods("DELETE")
}