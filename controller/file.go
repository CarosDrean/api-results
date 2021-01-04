package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileController struct {}

func (c FileController) DownloadZIPOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var filter models.Filter
	_ = json.NewDecoder(r.Body).Decode(&filter)
	res, err := db.ServiceDB{}.GetAllPatientsWithOrganizationFilter(filter)
	if err != nil {
		returnErr(w, err, "obtener todos organization filter")
		return
	}
	paths := make([]string, 0)
	for _, e := range res {
		petition := models.PetitionFile{
			Exam:        filter.Data,
			ServiceID:   e.ID,
			DNI:         e.DNI,
			NameComplet: e.FirstLastName + " " + e.SecondLastName + " " + e.Name,
			ServiceDate: e.ServiceDate,
		}
		path, err := c.assemblyFilePath(petition)
		if err == nil {
			paths = append(paths, path)
		}
	}
	output := "\\temp\\" + filter.Data + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	err = utils.ZipFiles(output, paths)
	if err != nil {
		returnErr(w, err, "crear zip")
		return
	}
	http.ServeFile(w, r, output)
}

func (c FileController) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var petition models.PetitionFile
	_ = json.NewDecoder(r.Body).Decode(&petition)

	filePath, err := c.assemblyFilePath(petition)
	if err != nil {
		returnErr(w, err, "obtener tarchivos")
		return
	}
	http.ServeFile(w, r, filePath)
}

func (c FileController) assemblyFilePath(petition models.PetitionFile) (string, error) {
	var nameFile string
	if strings.Contains(petition.Exam, "PRUEBA RAPIDA") {
		nameFile = constants.RoutePruebaRapida + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-PRUEBA-RAPIDA-" + constants.IdPruebaRapida + ".pdf"
	} else if strings.Contains(petition.Exam, "INTERCONSULTA"){
		nameFile = constants.RouteInterconsulta + petition.ServiceID + "-" + petition.NameComplet + ".pdf"
	} else if strings.Contains(petition.Exam, "INFORME MEDICO"){
		nameFile = constants.RouteInformeMedico + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "FMT2")
	} else if strings.Contains(petition.Exam, "CERTIFICADO SIN DX"){
		nameFile = constants.RouteCertificateSinDX + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "CAPSD")
	} else if strings.Contains(petition.Exam, "CERTIFICADO 312"){
		nameFile = constants.RouteCertificate312 + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "CAP")
	} else if strings.Contains(petition.Exam, "HISTORIA CLINICA") {
		nameFile = constants.RouteHistory + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "HISTORIA")
	}
	if len(nameFile) == 0 {
		return "", errors.New("no aceptado")
	}
	if _, err := os.Stat(nameFile); err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("no existe")
		}
	}

	return nameFile, nil
}

func (c FileController) assemblyFileNameExtra(idService string, dni string, parent string) string {
	fmt.Println(idService)
	person, _ := db.PersonDB{}.GetFromDNI(dni)
	service, _ := db.ServiceDB{}.Get(idService)
	protocol, _ := db.ProtocolDB{}.Get(service.ProtocolID)
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


