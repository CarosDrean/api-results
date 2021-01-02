package routes

import (
	statistic "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func statisticRoutes(s *mux.Router) {
	ctrl := statistic.StatisticController{}
	s.HandleFunc("/filter/", mid.CheckSecurity(ctrl.GetServiceDiseaseByProtocol)).Methods("POST")
}
