package controller

import (
	"github.com/CarosDrean/api-results.git/db"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
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

	organizationName := organization.Name
	personName := patients[0].Name
	date := services[0].ServiceDate
	dates := strings.Split(date, "T")
	log.Println(dates[0])
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dates[0])
	log.Println(t)

	namePDF := organizationName + " - " + personName + " - " + t.Format("02 Febrero, 2006") + ".pdf"
	log.Println(namePDF)
}
