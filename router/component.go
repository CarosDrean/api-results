package routes

import (
	component "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func componentRoutes(s *mux.Router) {
	s.HandleFunc("/all/{id}", mid.CheckSecurity(component.GetComponentsCategoryId)).Methods("GET")
}

