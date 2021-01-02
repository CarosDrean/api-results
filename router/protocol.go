package routes

import (
	protocol "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func protocolRoutes(s *mux.Router) {
	ctrl := protocol.ProtocolController{DB: db.ProtocolDB{}}
	s.HandleFunc("/all/{idLocation}", mid.CheckSecurity(ctrl.GetAllLocation)).Methods("GET")
	s.HandleFunc("/all-organization/{idOrganization}", mid.CheckSecurity(ctrl.GetAllOrganization)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(ctrl.Get)).Methods("GET")
}