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

func (c FileController) UploadAndSendZipOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var filter models.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		returnErr(w, err, "decoder filter")
		return
	}

	patientsOrganization, err := db.ServiceDB{}.GetAllPatientsWithOrganizationFilter(filter.ID, filter)
	if err != nil {
		returnErr(w, err, "obtener todos pacientes")
		return
	}

	filePaths := c.getFilePaths(patientsOrganization, filter.Data)
	if len(filePaths) == 0 {
		returnErr(w, nil, "sin elementos")
		return
	}

	fileName := filter.Data + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + ".zip"

	output := "temp\\" + fileName

	err = utils.CreateZip(output, filePaths)
	if err != nil {
		returnErr(w, err, "comprimir archivos")
		return
	}

	defer os.Remove(output)

	tokenUser := r.Header.Get("Authorization")

	responseFile, err := utils.UploadFile(constants.RouteUploadFile, output, tokenUser)
	if err != nil {
		returnErr(w, err, "subir archivo")
		return
	}

	mailFile := models.MailFile{
		Email:           filter.DataTwo,
		FilenameUpload:  responseFile.Data,
		Description:     "Recopilación de Historias Clinicas",
		NameFileSending: "Historias-Clinicas",
		FormatFile:      responseFile.Format,
	}

	mailResponse, err := c.sendZipOrganization(mailFile, tokenUser)
	if err != nil {
		returnErr(w, err, "enviar email")
		return
	}

	_ = json.NewEncoder(w).Encode(mailResponse)
}

func (c FileController) sendZipOrganization(mailFile models.MailFile, token string) (models.MailResponse, error){
	dataMailFile, err := json.Marshal(mailFile)
	if err != nil {
		return models.MailResponse{}, err
	}

	dataResponse, err := utils.SendMail(dataMailFile, constants.RouteSendFile, token)
	if err != nil {
		return models.MailResponse{}, err
	}

	mailResponse := models.MailResponse{}
	if err := mailResponse.Unmarshal(dataResponse); err != nil {
		return models.MailResponse{}, err
	}

	return mailResponse, nil
}

func (c FileController) getFilePaths(patients []models.ServicePatient, exam string) []string {
	paths := make([]string, 0)
	for _, patient := range patients {
		petition := models.PatientFile{
			Exam:        exam,
			ServiceID:   patient.ID,
			DNI:         patient.DNI,
			NameComplet: patient.FirstLastName + " " + patient.SecondLastName + " " + patient.Name,
			ServiceDate: patient.ServiceDate,
		}

		path, err := c.makeFilePath(petition)
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
	paths := c.getFilePaths(res, filter.Data)
	if len(paths) == 0 {
		log.Println("Sin elementos")
		return
	}
	output := "temp\\" + filter.Data + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	err = utils.CreateZip(output, paths)
	if err != nil {
		log.Println(fmt.Sprintf("Comprimir %s", err))
		return
	}
	http.ServeFile(w, r, output)
	_ = os.Remove(output)
}

func (c FileController) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var petition models.PatientFile
	_ = json.NewDecoder(r.Body).Decode(&petition)

	filePath, err := c.makeFilePath(petition)
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (c FileController) makeFilePath(patient models.PatientFile) (string, error) {
	var filePath string
	if strings.Contains(patient.Exam, "PRUEBA RAPIDA") {
		filePath = constants.RoutePruebaRapida + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA-RAPIDA-" + constants.IdPruebaRapida + ".pdf"
	} else if strings.Contains(patient.Exam, "INTERCONSULTA") {
		filePath = constants.RouteInterconsulta + patient.ServiceID + "-" + patient.NameComplet + ".pdf"
	} else if strings.Contains(patient.Exam, "INFORME MEDICO") {
		filePath = constants.RouteInformeMedico + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "FMT2")
	} else if strings.Contains(patient.Exam, "CERTIFICADO SIN DX") {
		filePath = constants.RouteCertificateSinDX + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "CAPSD")
	} else if strings.Contains(patient.Exam, "CERTIFICADO 312") {
		filePath = constants.RouteCertificate312 + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "CAP")
	} else if strings.Contains(patient.Exam, "HISTORIA CLINICA") {
		filePath = constants.RouteHistory + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "HISTORIA")
	}else if strings.Contains(patient.Exam, "PDF ADMINISTRATIVO") {
		filePath = constants.RoutePDFAdministrative + patient.DNI + " - " + formatDate(patient.ServiceDate) + ".pdf"
	}

	if strings.Contains((patient.Exam), "PRUEBA HISOPADO") {
		filePath = constants.RoutePruebaHisopado + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA-RAPIDA-HISOPADO-" + constants.IdPruebaHisopado + ".pdf"
	} else if strings.Contains(patient.Exam, "HOLOELECTRO") {
		filePath = constants.RouteCardio + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-SERVICIOS-" + constants.IdCardio + ".pdf"
	} else if strings.Contains(patient.Exam, "HOLTER") {
		filePath = constants.RouteHolter + c.assemblyFileDate(patient.ServiceID, patient.DNI, "HOLTER")
	} else if strings.Contains(patient.Exam, "ELECTROCARDIOGRAMA") {
		filePath = constants.RouteElectro + c.assemblyFileDate(patient.ServiceID, patient.DNI, "ELECTROCARDIOGRAMA")
	} else if strings.Contains(patient.Exam, "MAPA") {
		filePath = constants.RouteMapa + c.assemblyFileDate(patient.ServiceID, patient.DNI, "MAPA")
	} else if strings.Contains(patient.Exam, "ECOCARDIOGRAMA") {
		filePath = constants.RouteEcocardiograma + c.assemblyFileDate(patient.ServiceID, patient.DNI, "ECOCARDIOGRAMA")
	} else if strings.Contains(patient.Exam, "PRUEBA ESFUERZO") {
		filePath = constants.RoutePruebaEsfuerzo + c.assemblyFileDate(patient.ServiceID, patient.DNI, "PRUEBA ESFUERZO")
	} else if strings.Contains(patient.Exam, "RIESGO QUIRURGICO") {
		filePath = constants.RouteRiesgo + c.assemblyFileDate(patient.ServiceID, patient.DNI, "RIESGO QUIRURGICO")
	} else if strings.Contains(patient.Exam, "MANUAL DE HOLORESULTS - ADMINISTRADOR") {
		filePath = constants.RoutePDF + "MANUAL DE HOLORESULTS - ADMINISTRADOR" + ".pdf"
	} else if strings.Contains(patient.Exam, "MANUAL DE HOLORESULTS - MEDICO") {
		filePath = constants.RoutePDF + "MANUAL DE HOLORESULTS - MEDICO" + ".pdf"
	}

	if len(filePath) == 0 {
		return "", errors.New("no aceptado")
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			if strings.Contains((patient.Exam), "PRUEBA HISOPADO") {
				filePath = constants.RoutePruebaHisopado + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA-RAPIDA-HISOPADO-" + constants.IdPruebaHisopadoAux + ".pdf"
			} else {
				return "", errors.New("no existe")
			}

		}

	}
	return filePath, nil

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
	if parent == "PDF ADMINISTRATIVO" {
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
