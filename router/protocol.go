package routes

import (
	protocol "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func protocolRoutes(s *mux.Router) {
	s.HandleFunc("/all/{idSystemUser}", mid.CheckSecurity(protocol.GetProtocolsWidthSystemUser)).Methods("GET")
}
