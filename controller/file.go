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
		nameFile = constants.RouteInformeMedico + getFileNameInformeMedico(petition.ServiceID, petition.DNI)
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

func getFileNameInformeMedico(idService string, dni string) string {
	patients := db.GetPatientFromDNI(dni)
	services := db.GetService(idService)
	protocol := db.GetProtocol(services[0].ProtocolID)
	organization := db.GetOrganization(protocol.OrganizationID)
	log.Println(organization.Name)

	organizationName := organization.Name
	personName := patients[0].FirstLastName + " " + patients[0].SecondLastName + " " + patients[0].Name
	date := services[0].ServiceDate
	dates := strings.Split(date, "T")
	log.Println(dates[0])
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dates[0])
	log.Println(t)

	namePDF := organizationName + "-" + personName + "-FMT2-" + t.Format("02 febrero, 2006") + ".pdf"
	log.Println(namePDF)
	return namePDF
}
