
package routes

import (
	// mid "github.com/CarosDrean/api-holosalud/middleware"
	user "github.com/CarosDrean/api-results.git/controller"
	"github.com/gorilla/mux"
)

func userRoutes(s *mux.Router) {
	// s.HandleFunc("/", user.GetUsers).Methods("GET")
	s.HandleFunc("/{id}", user.GetUser).Methods("GET")
	s.HandleFunc("/", user.CreateUser).Methods("POST")
	// s.HandleFunc("/{id}", mid.CheckSecurity(user.UpdateUser)).Methods("PUT")
	// s.HandleFunc("/{id}", mid.CheckSecurity(user.DeleteUser)).Methods("DELETE")
}