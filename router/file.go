package routes

import (
	file "github.com/CarosDrean/api-results.git/controller"
	mid "github.com/CarosDrean/api-results.git/middleware"
	"github.com/gorilla/mux"
)

func fileRoutes(s *mux.Router) {
	ctrl := file.FileController{}
	s.HandleFunc("/", mid.CheckSecurity(ctrl.DownloadPDF)).Methods("POST")
	s.HandleFunc("/excel", mid.CheckSecurity(ctrl.DownloadExcelMatriz)).Methods("POST")
	s.HandleFunc("/all", mid.CheckSecurity(ctrl.DownloadZIPOrganization)).Methods("POST")
	s.HandleFunc("/send-file", mid.CheckSecurity(ctrl.UploadAndSendZipOrganization)).Methods("POST")
}
