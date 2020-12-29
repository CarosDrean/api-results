package controller

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileController struct {}

func (c FileController) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var petition models.PetitionFile
	_ = json.NewDecoder(r.Body).Decode(&petition)

	var nameFile string
	if strings.Contains(petition.Exam, "PRUEBA RAPIDA") {
		nameFile = constants.RoutePruebaRapida + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-PRUEBA-RAPIDA-" + constants.IdPruebaRapida + ".pdf"
	} else if strings.Contains(petition.Exam, "INTERCONSULTA"){
		nameFile = constants.RouteInterconsulta + petition.ServiceID + "-" + petition.NameComplet + ".pdf"
	} else if strings.Contains(petition.Exam, "INFORME MEDICO"){
		nameFile = constants.RouteInformeMedico + c.getFileNameReportMedicoAndCertificate(petition.ServiceID, petition.DNI, "FMT2")
	} else if strings.Contains(petition.Exam, "CERTIFICADO SIN DX"){
		nameFile = constants.RouteCertificateSinDX + c.getFileNameReportMedicoAndCertificate(petition.ServiceID, petition.DNI, "CAPSD")
	} else if strings.Contains(petition.Exam, "CERTIFICADO 312"){
		nameFile = constants.RouteCertificate312 + c.getFileNameReportMedicoAndCertificate(petition.ServiceID, petition.DNI, "CAP")
	} else if strings.Contains(petition.Exam, "HISTORIA CLINICA") {
		nameFile = constants.RouteHistory + c.getFileNameReportMedicoAndCertificate(petition.ServiceID, petition.DNI, "HISTORIA")
	}
	if len(nameFile) == 0 {
		log.Println("no aceptado")
		return
	}
	if _, err := os.Stat(nameFile); err != nil {
		if os.IsNotExist(err) {
			log.Println("no existe")
			return
		}
	}

	fp := path.Join(nameFile)
	http.ServeFile(w, r, fp)
}

func (c FileController) getFileNameReportMedicoAndCertificate(idService string, dni string, parent string) string {
	fmt.Println(idService)
	person, _ := db.PersonDB{}.GetFromDNI(dni)
	service, _ := db.ServiceDB{}.Get(idService)
	protocol := db.GetProtocol(service.ProtocolID)
	organization, _ := db.OrganizationDB{}.Get(protocol.OrganizationID)

	organizationName := organization.Name
	personName := person.FirstLastName + " " + person.SecondLastName + " " + person.Name

	date := service.ServiceDate
	dates := strings.Split(date, "T")
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dates[0])
	year, month, day := t.Date()
	td := strconv.Itoa(day) + " " + getMonth(month.String()) + ",  " + strconv.Itoa(year)

	namePDF := organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"
	if parent == "CAPSD" {
		personName = person.FirstLastName + " " + person.SecondLastName + ", " + person.Name
		namePDF = organizationName + " -" + personName + "-" + parent + "-" + td + ".pdf"
	}
	if  parent == "CAP" {
		personName = person.FirstLastName + " " + person.SecondLastName + ", " + person.Name
		namePDF = organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"
	}
	if  parent == "HISTORIA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		namePDF = organizationName + " - " + personName + " - " + td + ".pdf"
	}
	return namePDF
}


