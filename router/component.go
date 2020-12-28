package routes

import (
	component "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func componentRoutes(s *mux.Router) {
	ctrl := component.ComponentController{DB: db.ComponentDB{}}
	s.HandleFunc("/all/{id}", mid.CheckSecurity(ctrl.GetAllCategoryId)).Methods("GET")
}

