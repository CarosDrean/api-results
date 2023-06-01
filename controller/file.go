package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileController struct{
	DB db.FileDB
}

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
		//filePath = constants.RouteHistory + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "HISTORIA")
		filePath = constants.RouteReportesMedicos + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "HISTORIA AUDITADA") {
		filePath = constants.RouteHistoryAudity + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "HISTORIA AUDITADA")
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

	if strings.Contains(patient.Exam, "PRUEBA RAPIDA - PARTICULAR") {
		filePath = constants.RoutePruebaRapidaParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA-RAPIDA-N007-ME000000491.pdf"
	} else if strings.Contains(patient.Exam, "PRUEBA HISOPADO - PARTICULAR") {
		filePath = constants.RouteInterconsultaParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA-RAPIDA-N009-ME000000529.pdf"
	} else if strings.Contains(patient.Exam, "INTERCONSULTA - PARTICULAR") {
		filePath = constants.RouteInterconsultaParticular + patient.ServiceID + "-" + patient.NameComplet + ".pdf"
	} else if strings.Contains(patient.Exam, "CONSULTA CARDIOLOGICA - PARTICULAR") {
		filePath = constants.RouteConsultaCardioParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-CONSULTA-CARDIOLOGICA-N009-ME000000534.pdf"
	} else if strings.Contains(patient.Exam, "HISTORIA CLINICA - PARTICULAR") {
		filePath = constants.RouteHistoryParticular + c.assemblyFileNameExtra(patient.ServiceID, patient.DNI, "HISTORIA")
	} else if strings.Contains(patient.Exam, "ELECTROCARDIOGRAMA - PARTICULAR") {
		filePath = constants.RouteElectroParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-ELECTROCARDIOGRAMA-" + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "HOLTER - PARTICULAR") {
		filePath = constants.RouteHolterParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-HOLTER-" + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "ECOCARDIOGRAMA - PARTICULAR") {
		filePath = constants.RouteEcocardiogramaParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-ECOCARDIOGRAMA-" + patient.ServiceID + ".pdf"
	}  else if strings.Contains(patient.Exam, "MAPA - PARTICULAR") {
		filePath = constants.RouteMapaParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-MAPA-" + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "PRUEBA ESFUERZO - PARTICULAR") {
		filePath = constants.RoutePruebaEsfuerzoParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-PRUEBA ESFUERZO-" + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "RIESGO QUIRURGICO - PARTICULAR") {
		filePath = constants.RouteRiesgoParticular + patient.DNI + "-" + formatDate(patient.ServiceDate) + "-RIESGO QUIRURGICO-" + patient.ServiceID + ".pdf"
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
	if parent == "HISTORIA AUDITADA" {
		personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name + " - " + strconv.Itoa(protocol.EsoType) + " - " + protocol.GroupOccupationId
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

func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

func (c FileController) DownloadExcelMatriz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var exs models.ExcelPetitionMatrizFile
	_ = json.NewDecoder(r.Body).Decode(&exs)

	filePath, err := c.ExcelMatriz(exs)
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeFile(w, r, filePath)
}


func (c FileController) ExcelMatriz(exs models.ExcelPetitionMatrizFile) (string, error){
	f := excelize.NewFile()

	streamWriter, err := f.NewStreamWriter("Sheet1")

	if err != nil {
		fmt.Println(err)
		//return
	}

	styleID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		fmt.Println(err)
		//return
	}

	z := 1
	a := [24]float64{7, 25, 15, 50, 9, 25, 150, 37, 45, 30, 50, 20, 4, 57, 29, 4, 57, 29, 4, 57, 29, 28, 28, 28}

	for _, p := range a {

		if err := streamWriter.SetColWidth(z, z, p); err != nil {
			fmt.Println(err)
			//return
		}

		z = z + 1
	}

	header := []interface{}{}

	for _, cell := range []string{
		"N°", "CODIGO DE ATENCION", "DNI", "APELLIDOS Y NOMBRES", "EDAD", "TIPO DE EXAMEN", "PROTOCOLO", "FECHA DE EVALUACION OCUPACIONAL",
		"PUESTO EN LA EMPRESA", "APTITUD LABORAL PARA EL PUESTO DE TRABAJO", " ", " ", "I1", "RECOMENDACION", "ESPECIALIDAD", "I2", "RECOMENDACION", "ESPECIALIDAD", "I3",
		"RECOMENDACION", "ESPECIALIDAD", "APTITUD ESPECIFICA", " ", "OBSERVACION GENERAL",
	} {
		header = append(header, cell)
	}

	if err := streamWriter.SetRow("A2", header, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	header2 := []interface{}{}

	for _, cell2 := range []string{
		"APTITUD MEDICA", "RESTRICCION", "CERT. ETAS", " ", " ", " ", " ", " ", " ", " ", " ", " ", "ALTURA", "ESPACIOS CONFINADOS",
	} {
		header2 = append(header2, cell2)
	}

	if err := streamWriter.SetRow("J3", header2, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	merges := map[string]string{
		"A2": "A3", "B2": "B3", "C2": "C3", "D2": "D3", "E2": "E3", "F2": "F3", "G2": "G3",
		"H2": "H3", "I2": "I3", "J2": "L2", "M2": "M3", "N2": "N3", "O2": "O3", "P2": "P3",
		"Q2": "Q3", "R2": "R3", "S2": "S3", "T2": "T3", "U2": "U3", "V2": "W2", "X2": "X3",
	}

	for x, v := range merges {
		if err := streamWriter.MergeCell(x, v); err != nil {
			fmt.Println(err)
			//return
		}
	}

	//items, _ := c.DB.GetMatrizOnline("2023-05-03", "2023-05-04")
	items, _ := c.DB.GetMatrizOnline(exs.Ini, exs.Fin, exs.OrganizationID)

	y := 4
	x := 1
	p := 0
	cont := 0

	RowCellValue := make([]interface{}, 0)

	for _, e := range items {
		ss := "A" + strconv.Itoa(y)
		ages, _ := time.Parse("2006-01-02", e.Bithdate[0:10])

		//---PROCESO-INTERCONSULTAS---
		ex, _ := c.DB.GetInterconsultas(e.VServiceid)

		arr1 := [3]string{}
		arr2 := [3]string{}

		for _, g := range ex {
			if strings.Contains(g.InterconsultaName, "I/C") == true {
				arr2[p] = g.InterconsultaName

				//---PROCESO-INTERCONSUTLAS-RECOMENDACIONES---
				exr, _ := c.DB.GetRecomendaciones(g.RepositorioDxId, e.VServiceid)
				reco := "---"

				for _, j := range exr {
					if len(j.RecomendationName) > 0 {
						reco = j.RecomendationName + "," + reco
					}
				}

				recoFilter1 := strings.Replace(reco, ",---", "", -1)
				arr1[p] = recoFilter1
				//----------------------------

			}

			p = p + 1

			if p == 2 {
				p = 0
			}
		}

		if len(arr2[0]) < 1 {
			arr2[0] = "---"
		}
		if len(arr2[1]) < 1 {
			arr2[1] = "---"
		}
		if len(arr2[2]) < 1 {
			arr2[2] = "---"
		}

		if len(arr1[0]) < 1 {
			arr1[0] = "---"
		}
		if len(arr1[1]) < 1 {
			arr1[1] = "---"
		}
		if len(arr1[2]) < 1 {
			arr1[2] = "---"
		}
		//----------------------------

		//---PROCESO-RESTRICCIONES---

		rest, _ := c.DB.GetRestricciones(e.VServiceid)
		acu := "---"

		for _, h := range rest {
			if h.RestrictionName != "." {
				acu = h.RestrictionName + "," + acu
			} else {

			}
		}

		acuFilter1 := strings.Replace(acu, ",---", "", -1)

		//----------------------------

		//---PROCESO-ALTURA-APTITUD---

		al, _ := c.DB.GetAlturaAptitud(e.VServiceid)
		alti := "---"

		for _, k := range al {
			if len(k.AptitudName) > 0 {
				alti = k.AptitudName
			} else {
				alti = "---"
			}
		}

		//----------------------------

		//---PROCESO-ESPACIOS-CONFINADOS-APTITUD---

		ec, _ := c.DB.GetAptitudEspaciosConfi(e.VServiceid)
		eco := "---"

		for _, cc := range ec {
			if len(cc.AptitudName) > 0 {
				if cc.AptitudName == "1" {
					eco = "APTO"
				}else{
					eco = "NO APTO"
				}
			} else {
				eco = "---"
			}
		}

		//----------------------------

		RowCellValue = append(make([]interface{}, 0), strconv.Itoa(x), e.VServiceid, e.DocNumber, e.PersonName, strconv.Itoa(Age(ages)), e.EsoName, e.ProtocolName, e.ServiceDate[0:10], e.PersonOcupation, e.Aptitude, acuFilter1, "NO APLICA", "1", arr1[0], arr2[0], "1", arr1[1], arr2[1], "1", arr1[2], arr2[2], alti, eco, "EVALUADO")

		if err := streamWriter.SetRow(ss, RowCellValue, excelize.RowOpts{StyleID: styleID}); err != nil {
			fmt.Println(err)
			//return
		}

		y = y + 1
		x = x + 1
		cont = cont + 3
	}

	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
		//return
	}

	if err := f.SaveAs("\\\\HOLO-SERVIDOR\\archivos sistema_2\\TEMPORAL\\" + exs.OrganizationID + ".xlsx"); err != nil {
		println(err.Error())
	}

	var filePath string

	filePath = "\\\\HOLO-SERVIDOR\\archivos sistema_2\\TEMPORAL\\" + exs.OrganizationID + ".xlsx";

	if len(filePath) == 0 {
		return "", errors.New("no aceptado")
	}

	return filePath, nil
}