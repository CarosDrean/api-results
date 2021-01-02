package routes

import (
	protocol "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func protocolSystemUserRoutes(s *mux.Router) {
	ctrl := protocol.ProtocolSystemUserController{}
	s.HandleFunc("/all/{idSystemUser}", mid.CheckSecurity(ctrl.GetSystemUser)).Methods("GET")
}
