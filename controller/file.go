package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileController struct{}

func (c FileController) SendZipOrganizationData(mailFile models.MailFile) error {
	data, _ := json.Marshal(mailFile)
	err := utils.SendMail(data, constants.RouteSendFile)
	if err != nil {
		return err
	}
	return nil
}

func (c FileController) SendZipOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var filter models.Filter
	_ = json.NewDecoder(r.Body).Decode(&filter)
	res, err := db.ServiceDB{}.GetAllPatientsWithOrganizationFilter(filter.ID, filter)
	if err != nil {
		log.Println(err)
		returnErr(w, err, "obtener todos pacientes")
		return
	}
	paths := c.GetPaths(res, filter.Data)
	if len(paths) == 0 {
		log.Println("Sin elementos")
		returnErr(w, err, "sin elementos")
		return
	}
	fileName := filter.Data + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + ".zip"
	output := "temp\\" + fileName
	err = utils.ZipFiles(output, paths)
	if err != nil {
		log.Println(fmt.Sprintf("Comprimir %s", err))
		returnErr(w, err, "comprimir archivos")
		return
	}
	err = utils.SendFileMail(filter.DataTwo, constants.RouteUploadFile, output)
	if err != nil {
		log.Println(err)
		returnErr(w, err, "subir archivo")
	}

	organization, _ := db.OrganizationDB{}.Get(filter.ID)
	mailFile := models.MailFile{
		From:     filter.DataTwo,
		File:     fileName,
		Business: organization.Name,
		DateFrom: filter.DateFrom,
		DateTo:   filter.DateTo,
	}
	err = c.SendZipOrganizationData(mailFile)
	if err != nil {
		log.Println(err)
		returnErr(w, err, "enviar email")
	}

	_ = os.Remove(output)
	_ = json.NewEncoder(w).Encode("enviado!")
}

func (c FileController) GetPaths(res []models.ServicePatient, exam string) []string {
	paths := make([]string, 0)
	for _, e := range res {
		petition := models.PetitionFile{
			Exam:        exam,
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
	return paths
}

func (c FileController) DownloadZIPOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var filter models.Filter
	_ = json.NewDecoder(r.Body).Decode(&filter)
	var err error
	res := make([]models.ServicePatient, 0)
	if filter.ID == "all" {
		res, err = db.ServiceDB{}.GetAllPatientsWithOrganizationFilter(filter.DataTwo, filter)
	} else {
		res, err = db.ServiceDB{}.GetAllPatientsWithProtocolFilter(filter.ID, filter, false)
	}

	if err != nil {
		log.Println(err)
		return
	}
	paths := c.GetPaths(res, filter.Data)
	if len(paths) == 0 {
		log.Println("Sin elementos")
		return
	}
	output := "temp\\" + filter.Data + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	err = utils.ZipFiles(output, paths)
	if err != nil {
		log.Println(fmt.Sprintf("Comprimir %s", err))
		return
	}
	http.ServeFile(w, r, output)
	_ = os.Remove(output)
}

func (c FileController) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var petition models.PetitionFile
	_ = json.NewDecoder(r.Body).Decode(&petition)

	filePath, err := c.assemblyFilePath(petition)
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (c FileController) assemblyFilePath(petition models.PetitionFile) (string, error) {
	var nameFile string
	if strings.Contains(petition.Exam, "PRUEBA RAPIDA") {
		nameFile = constants.RoutePruebaRapida + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-PRUEBA-RAPIDA-" + constants.IdPruebaRapida + ".pdf"
	} else if strings.Contains(petition.Exam, "INTERCONSULTA") {
		nameFile = constants.RouteInterconsulta + petition.ServiceID + "-" + petition.NameComplet + ".pdf"
	} else if strings.Contains(petition.Exam, "INFORME MEDICO") {
		nameFile = constants.RouteInformeMedico + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "FMT2")
	} else if strings.Contains(petition.Exam, "CERTIFICADO SIN DX") {
		nameFile = constants.RouteCertificateSinDX + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "CAPSD")
	} else if strings.Contains(petition.Exam, "CERTIFICADO 312") {
		nameFile = constants.RouteCertificate312 + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "CAP")
	} else if strings.Contains(petition.Exam, "HISTORIA CLINICA") {
		nameFile = constants.RouteHistory + c.assemblyFileNameExtra(petition.ServiceID, petition.DNI, "HISTORIA")
	} else if strings.Contains(petition.Exam, "PRUEBA HISOPADO") {
		nameFile = constants.RoutePruebaHisopado + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-PRUEBA-RAPIDA-HISOPADO-" + constants.IdPruebaHisopado + ".pdf"
	} else if strings.Contains(petition.Exam, "HOLOELECTRO") {
		nameFile = constants.RouteCardio + petition.DNI + "-" + formatDate(petition.ServiceDate) + "-SERVICIOS-" + constants.IdCardio + ".pdf"
	} else if strings.Contains(petition.Exam, "HOLTER") {
		nameFile = constants.RouteHolter + c.assemblyFileDate(petition.ServiceID, petition.DNI, "HOLTER")
	} else if strings.Contains(petition.Exam, "ELECTROCARDIOGRAMA") {
		nameFile = constants.RouteElectro + c.assemblyFileDate(petition.ServiceID, petition.DNI, "ELECTROCARDIOGRAMA")
	} else if strings.Contains(petition.Exam, "MAPA") {
		nameFile = constants.RouteMapa + c.assemblyFileDate(petition.ServiceID, petition.DNI, "MAPA")
	} else if strings.Contains(petition.Exam, "ECOCARDIOGRAMA") {
		nameFile = constants.RouteEcocardiograma + c.assemblyFileDate(petition.ServiceID, petition.DNI, "ECOCARDIOGRAMA")
	} else if strings.Contains(petition.Exam, "PRUEBA ESFUERZO") {
		nameFile = constants.RoutePruebaEsfuerzo + c.assemblyFileDate(petition.ServiceID, petition.DNI, "PRUEBA ESFUERZO")
	}else if strings.Contains(petition.Exam, "RIESGO QUIRURGICO") {
		nameFile = constants.RouteRiesgo + c.assemblyFileDate(petition.ServiceID, petition.DNI, "RIESGO QUIRURGICO")

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
	dayString := strconv.Itoa(day)
	if len(dayString) == 1 {
		dayString = "0" + strconv.Itoa(day)
	}
	td := dayString + " " + getMonth(month.String()) + ",  " + strconv.Itoa(year)

	namePDF := organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"

	if parent == "CAPSD" {
		personName = person.FirstLastName + " " + person.SecondLastName + ", " + person.Name
		namePDF = organizationName + " -" + personName + "-" + parent + "-" + td + ".pdf"
	}
	if parent == "CAP" {
		personName = person.FirstLastName + " " + person.SecondLastName + ", " + person.Name
		namePDF = organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"
	}
	if parent == "HISTORIA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		namePDF = organizationName + " - " + personName + " - " + td + ".pdf"
	}

	return namePDF
}
func (c FileController) assemblyFileDate(idService string, dni string, parent string) string {
	person, _ := db.PersonDB{}.GetFromDNI(dni)
	service, _ := db.ServiceDB{}.Get(idService)
	protocol, _ := db.ProtocolDB{}.Get(service.ProtocolID)
	organization, _ := db.OrganizationDB{}.Get(protocol.OrganizationID)

	organizationName := organization.Name
	personName := person.FirstLastName + " " + person.SecondLastName + " " + person.Name
	personDoc := person.DNI

	date := service.ServiceDate
	dates := strings.Split(date, "T")
	layout := "2006-01-02"
	t, _ := time.Parse(layout, dates[0])
	year, month, day := t.Date()
	dayString := strconv.Itoa(day)

	td := dayString + "" + strconv.Itoa(int(month)) + "" + strconv.Itoa(year)

	namePDF := organizationName + "-" + personName + "-" + parent + "-" + td + ".pdf"


	if parent == "HOLTER" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"
	}
	if parent == "ELECTROCARDIOGRAMA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"
	}
	if parent == "MAPA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"
	}
	if parent == "ECOCARDIOGRAMA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"
	}
	if parent == "PRUEBA ESFUERZO" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"
	}
	if parent == "RIESGO QUIRURGICO" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name
		personDoc = person.DNI
		namePDF = personDoc + "-" + td + "-" + parent + "-" + service.ProtocolID + ".pdf"

	}
	fmt.Println(namePDF)
	return namePDF
}
