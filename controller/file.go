package controller

import (
	"github.com/CarosDrean/api-results.git/db"
	"log"
	"net/http"
	"path"
)

func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("\\\\DESKTOP-QD7QM2Q\\archivos sistema_2\\Consolidado\\PRUEBA RAPIDA\\CONTRATISTA LOS MAGNIFICOS S.A.C. - CUNYAS VILA JEAN CARLOS - 03 septiembre,  2020.pdf")
	http.ServeFile(w, r, fp)
}

func GetData(dni string){
	patients := db.GetPatientFromDNI(dni)
	services := db.GetServiceWidthPersonID(patients[0].ID)
	protocol := db.GetProtocol(services[0].ProtocolID)
	organization := db.GetOrganization(protocol.OrganizationID)
	log.Println(organization.Name)
}
