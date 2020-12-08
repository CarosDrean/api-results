package controller

import (
	"encoding/json"
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

func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var petition models.PetitionFile
	_ = json.NewDecoder(r.Body).Decode(&petition)

	var nameFile string
	if strings.Contains(petition.Exam, "PRUEBA RAPIDA") {
		nameFile = constants.RoutePruebaRapida + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-PRUEBA-RAPIDA-" + constants.IdPruebaRapida + ".pdf"
	} else if strings.Contains(petition.Exam, "INTERCONSULTA"){
		nameFile = constants.RouteInterconsulta + petition.ServiceID + "-" + petition.NameComplet + ".pdf"
	} else if strings.Contains(petition.Exam, "INFORME MEDICO"){
		nameFile = constants.RouteInformeMedico + getFileNameInformeMedicoAndCertificate(petition.ServiceID, petition.DNI, "FMT2")
	} else if strings.Contains(petition.Exam, "CERTIFICADO SIN DX"){
		nameFile = constants.RouteCertificateSinDX + getFileNameInformeMedicoAndCertificate(petition.ServiceID, petition.DNI, "CAPSD")
	} else if strings.Contains(petition.Exam, "CERTIFICADO 312"){
		nameFile = constants.RouteCertificate312 + getFileNameInformeMedicoAndCertificate(petition.ServiceID, petition.DNI, "CAP")
	} else if strings.Contains(petition.Exam, "HISTORIA CLINICA") {

	}
	log.Println(nameFile)
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
	// fp := path.Join("\\\\DESKTOP-QD7QM2Q\\archivos sistema_2\\Consolidado\\PRUEBA RAPIDA\\CONTRATISTA LOS MAGNIFICOS S.A.C. - CUNYAS VILA JEAN CARLOS - 03 septiembre,  2020.pdf")
	http.ServeFile(w, r, fp)
}

func formatDate(date string) string {
	data := strings.FieldsFunc(date, SplitTwo)
	return data[2]+data[1]+data[0]
}

func SplitTwo(r rune) bool {
	return r == '-' || r == 'T'
}

func getFileNameInformeMedicoAndCertificate(idService string, dni string, parent string) string {
	patients := db.GetPatientFromDNI(dni)
	services := db.GetService(idService, db.NQGetService)
	protocol := db.GetProtocol(services[0].ProtocolID)
	organization := db.GetOrganization(protocol.OrganizationID)

	organizationName := organization.Name
	personName := patients[0].FirstLastName + " " + patients[0].SecondLastName + " " + patients[0].Name

	date := services[0].ServiceDate
	dates := strings.Split(date, "T")
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dates[0])
	year, month, day := t.Date()
	td := strconv.Itoa(day) + " " + getMonth(month.String()) + ",  " + strconv.Itoa(year)

	namePDF := organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"
	if parent == "CAPSD" {
		personName = patients[0].FirstLastName + " " + patients[0].SecondLastName + ", " + patients[0].Name
		namePDF = organizationName + " -" + personName + "-" + parent + "-" + td + ".pdf"
	}
	if  parent == "CAP" {
		personName = patients[0].FirstLastName + " " + patients[0].SecondLastName + ", " + patients[0].Name
		namePDF = organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"
	}
	log.Println(namePDF)
	return namePDF
}

func getMonth(month string) string {
	switch month {
	case "January":
		return "enero"
	case "February":
		return "febrero"
	case "March":
		return "marzo"
	case "April":
		return "abril"
	case "May":
		return "mayo"
	case "June":
		return "junio"
	case "July":
		return "julio"
	case "August":
		return "agosto"
	case "September":
		return "septiembre"
	case "October":
		return "octubre"
	case "November":
		return "noviembre"
	case "December":
		return "diciembre"
	default:
		return month
	}
}
