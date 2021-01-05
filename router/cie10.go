package routes

import (
	cie10 "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func cie10Routes(s *mux.Router) {
	ctrl := cie10.CIE10Controller{DB: db.CIE10DB{}}
	s.HandleFunc("/", mid.CheckSecurity(ctrl.GetAll)).Methods("GET")
}
