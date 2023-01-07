package routes

import (
	petition "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func petitionRoutes(s *mux.Router) {
	ctrl := petition.PetitionController{DB: db.PetitionDB{}}
	s.HandleFunc("/", mid.RoleInternalAdminOrTemp(ctrl.Create)).Methods("POST")
}