package routes

import (
	protocol "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func protocolRoutes(s *mux.Router) {
	s.HandleFunc("/all/{idLocation}", mid.CheckSecurity(protocol.GetProtocolsWidthLocation)).Methods("GET")
	s.HandleFunc("/all-organization/{idOrganization}", mid.CheckSecurity(protocol.GetProtocolsWidthOrganization)).Methods("GET")
	s.HandleFunc("/{id}", mid.CheckSecurity(protocol.GetProtocol)).Methods("GET")
}