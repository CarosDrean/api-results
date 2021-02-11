package routes

import (
	service "github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/db"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func serviceRoutes(s *mux.Router) {
	ctrl := service.ServiceController{DB: db.ServiceDB{}}
	s.HandleFunc("/all/{idProtocol}", mid.CheckSecurity(ctrl.GetAllPatientsWithProtocol)).Methods("GET")
	s.HandleFunc("/all-covid/{docNumber}", mid.CheckSecurity(ctrl.GetAllCovid)).Methods("GET")
	s.HandleFunc("/all-organization/{id}", mid.CheckSecurity(ctrl.GetAllPatientsWithOrganization)).Methods("GET")
	s.HandleFunc("/all-organization/", mid.CheckSecurity(ctrl.GetAllPatientsWithOrganizationFilter)).Methods("POST")
	s.HandleFunc("/filter/", mid.CheckSecurity(ctrl.GetAllPatientsWithProtocolFilter)).Methods("POST")
	s.HandleFunc("/filter-date/", mid.RoleInternalAdmin(ctrl.GetAllDiseaseFilterDate)).Methods("POST")
	s.HandleFunc("/all-date/", mid.CheckSecurity(ctrl.GetAllDate)).Methods("POST")
	s.HandleFunc("/all-exam-detail/", mid.CheckSecurity(ctrl.GetAllExamsDetail)).Methods("POST")
}
