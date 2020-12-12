package routes

import (
	statistic "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func statisticRoutes(s *mux.Router) {
	s.HandleFunc("/filter/", mid.CheckSecurity(statistic.GetStatisticsServiceDiseaseByProtocol)).Methods("POST")
}
