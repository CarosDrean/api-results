package routes

import (
	patient "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func personRoutes(s *mux.Router) {
	ctrl := patient.PersonController{DB: db.PersonDB{}}
	s.HandleFunc("/all/{idProtocol}", mid.CheckSecurity(ctrl.GetAllProtocol)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Get)).Methods("GET")
	s.HandleFunc("/dni/{dni}", mid.CheckSecurity(ctrl.GetFromDNI)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.UpdatePassword)).Methods("PUT")
}
