package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"github.com/xuri/excelize/v2"
)

type FileController struct {
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

func (c FileController) sendZipOrganization(mailFile models.MailFile, token string) (models.MailResponse, error) {
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
		//filePath = constants.RouteHistoryAudity + patient.ServiceID + ".pdf"
	} else if strings.Contains(patient.Exam, "PDF ADMINISTRATIVO") {
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
	} else if strings.Contains(patient.Exam, "MAPA - PARTICULAR") {
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
		//personName = person.FirstLastName + " " + person.SecondLastName + " " + person.Name + " - " + strconv.Itoa(protocol.EsoType) + " - " + protocol.GroupOccupationId
		namePDF = organizationName + "\\" + idService + ".pdf"
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

func (c FileController) DownloadExcelMatrizGrande(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var exs models.ExcelPetitionMatrizFile
	_ = json.NewDecoder(r.Body).Decode(&exs)

	filePath, err := c.ExcelMatrizGrande(exs)
	if err != nil {
		log.Println(err)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (c FileController) ExcelMatriz(exs models.ExcelPetitionMatrizFile) (string, error) {
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
				} else {
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

	filePath = "\\\\HOLO-SERVIDOR\\archivos sistema_2\\TEMPORAL\\" + exs.OrganizationID + ".xlsx"

	if len(filePath) == 0 {
		return "", errors.New("no aceptado")
	}

	return filePath, nil
}

func (c FileController) ExcelMatrizGrande(exs models.ExcelPetitionMatrizFile) (string, error) {
	f := excelize.NewFile()

	streamWriter, err := f.NewStreamWriter("Sheet1")

	if err != nil {
		fmt.Println(err)
		//return
	}

	styleID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	if err != nil {
		fmt.Println(err)
		//return
	}

	z := 1
	a := [62]float64{4.29, 7.57, 8.57, 8.86, 7.86, 9, 10.86, 4.43, 8.57, 14.29, 10.29, 17.14, 10.57,
		7, 10.71, 8.57, 9, 10.57, 11, 11, 35.14, 14.71, 14.29, 12.86, 11, 11, 11,
		11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
		9.29, 9.43, 9.43, 5.14}

	for _, p := range a {

		if err := streamWriter.SetColWidth(z, z, p); err != nil {
			fmt.Println(err)
			//return
		}

		z = z + 1
	}

	header := []interface{}{}

	for _, cell := range []string{
		"DATOS GENERALES", "", "", "", "", "", "", "", "", "", "", "", "",
		"RESULTADOS VAMO", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"ODONTOLOGIA", "AUDIOMETRIA", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "RAYOS X TORAX", "", "", "ESPIROMETRIA", "", "", "",
		"CARDIOLOGIA", "", "", "", "", "", "", "", "OFTALMOLOGIA", "", "", "", "", "", "", "", "", "", "", "", "", "", "EVALUACIÓN CLINICA", "", "", "", "", "", "", "",
		"", "", "INMUNIZACIONES", "", "", "", "", "COVID", "", "", "", "", "Descarte TBC", "", "PSICOLOGIA", "", "PERFIL CONDUCTOR", "", "", "", "", "", "",
		"TRABAJOS EN ALTURA ESTRUCTURAL > 1.8 metros", "", "", "", "", "", "EVALUACION NEUROLOGICA POR MEDICO NEUROLOGO", "EEG", "MANIPULADOR DE ALIMENTOS Y RESIDUOS (Apto / No apto)",
		"Llenado por Centro Médico Autorizado (CMA) por HB", "", "Llenado por Médico Ocupacional (MO) de UMC o de Empresa Contratista", "", "", "",
		"METALES PESADOS (si no aplica colocar NA)", "", "", "", "SATISFACCION DEL USUARIO (valor numérico absoluto)", "PROYECTO DE TRABAJADOR", "CLINICA ORIGEN", "MANIPULADOR DE ALIMENTOS",
		"", "", "", "", "HISTORICO OCUPACIONAL", "", "", "EVAL. PARA TRABAJOS EN ESPACIOS CONFINADOS", "EVAL ANTROPOLÓGICA", "", "TEST DE EPWORTH", "APTITUD OSTEOMUSCULAR",
		"DERMATOLOGICO", "OTOSCOPIA (EXAMEN ORL)", "", "", "", "TEST DE ESTRÉS SQR",
	} {
		header = append(header, cell)
	}

	if err := streamWriter.SetRow("A1", header, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	header2 := []interface{}{}

	for _, cell2 := range []string{
		"", "Nro DNI/CE/PAS", "Apellido materno", "Apellido paterno", "Primer nombre", "Segundo nombre", "Fecha de nacimiento (dd/mm/aaaa)",
		"Edad\n(años)", "Género", "Departamento de procedencia", "Direccion domiciliaria actualizada",
		"Empresa", "Ocupación / Cargo",

		"Tipo de EMO (EMPO, EMOA, EMOR, OTROS)", "Fecha de evaluación de EMO (dd/mm/aaaa)",
		"Resultado de Aptitud Médica", "Fecha de vencimiento (dd/mm/aaaa)", "Observaciones a levantar en 01 mes",
		"Observaciones a levantar en 03 meses", "Observaciones a levantar en 06 meses", "Resultado de Aptitud Específica (si es que aplica)",
		"Restricción 01", "Restricción 02", "Restricción 03", "Restricción 04", "Restricción 05", "Restricción 06",

		"HEMOGRAMA", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"Grupo sanguineo y Factor RH", "Glucosa (mg/dl)", "Hb glicosilada HbA1c (%)", "RPR",

		"EXAMEN DE ORINA", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "Prueba de HCGB (Edad Fertil  - Mujeres) mIU/mL", "Perfil lipídico", "", "", "", "Nº piezas con caries", "Oido derecho",
		"", "", "", "", "", "", "", "", "", "Oido izquierdo", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "Agudeza visual Cerca", "", "", "", "Agudeza visual Lejos", "", "", "", "Agudeza binocular", "", "Visión de colores (Test de Ishihara)",
		"Visión esteroscópica (segundos)", "Campos visuales", "", "Nutricional", "", "", "", "Antecedentes", "", "", "", "", "", "Tetanos", "Fiebre Tifoidea", "Hepatitis A",
		"Hepatitis B: 3 dosis", "Influenza: anual", "1ra Dosis", "2da Dosis", "3ra Dosis", "4ta Dosis", "Refuerzo", "Presenta cuadro Sintomático Respiratorio",
		"Baciloscopía (2 muestras, zonas endémicas / criterio médico)", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Nombre del CMA donde se realizó el EMO",
		"Médico Responsable de la evaluación del EMO", "FECHA DE REVISION MO CONSTANCIA / MO CONTRATA", "", "LEVANTAMIENTO DE OBSERVACION MO CONSTANCIA / MO CONTRATA", "",
	} {
		header2 = append(header2, cell2)
	}

	if err := streamWriter.SetRow("A2", header2, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	header4 := []interface{}{}

	for _, cell4 := range []string{
		"COMPLETO", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"Toxicológico", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "Fasciograma", "", "", "", "Alergias a medicamentos",
	} {
		header4 = append(header4, cell4)
	}

	if err := streamWriter.SetRow("BK3", header4, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	header3 := []interface{}{}

	for _, cell3 := range []string{
		"LEUCOCITOS - Cel/uL", "HEMATIES - Cel/uL", "HEMOGLOBINA - g/dL", "HEMATOCRITO - %", "NEUTROFILOS SEGMENTADOS - %", "NEUTROFILOS SEGMENTADOS - Cel/mL",
		"PLAQUETAS - Cel/uL", "VCM - fL", "HCM - pg", "CCMH - g/dL", "RDW - %", "VPM - fL", "BASTONES - %", "LINFOCITOS - %", "MONOCITOS - %", "EOSINOFILOS - %",
		"BASOFILOS - %", "METAMIELOCITOS - %", "MIELOCITOS - %", "PROMIELOCITOS - %", "BLASTOS - %", "BASTONES - Cel/mL", "LINFOCITOS - Cel/mL", "MONOCITOS - Cel/mL",
		"EOSINOFILOS - Cel/mL", "BASOFILOS - Cel/mL", "MIELOCITOS - Cel/mL", "METAMIELOCITOS - Cel/mL", "PROMIELOCITOS - Cel/mL", "BLASTOS - Cel/mL",
		"COMPROBACION RECUENTO 100%", "", "", "", "", "Color", "Aspecto", "Densidad", "Ph", "Glucosa", "Bilirrubina", "Cuerpos Cetonicos", "Proteinas", "Urobilinogeno",
		"Nitritos", "HEMOGLOBINA", "LEUCOCITOS - Cel/Cam", "HEMATIES - Hemat/mL", "CELULAS EPITELIALES", "LEVADURAS", "CRISTALES", "C.AC-URICO", "C.FOSF-AMORFOS", "C.URAT-AMORFOS",
		"C.OX-CALCIO", "C.FOSF-TRIPLES", "CILINDROS", "HIALINOS - /CAMPO", "GRANULOSOS - /CAMPO", "FILAMENTO MUCOIDES", "GERMENES", "Cocaina", "Marihuana",
		"Anfetamina", "Metanfetamina", "Metadora", "Morfina", "Feciclidina", "Barbituricos", "Benzodiacepinas", "Antidepresivos", "", "Colesterol Total (mg/dl)",
		"HDL-Col (mg/dl)", "LDL-Col (mg/dl)", "Trigliceridos (mg/dl)", "", "500", "1K", "2K", "3K", "4K", "6K", "8K", "STS", "Interpretación Clínica", "Interpretación Ocupacional",
		"500", "1K", "2K", "3K", "4K", "6K", "8K", "STS", "Interpretación Clínica", "Interpretación Ocupacional", "OIT", "Hallazgos", "Dx", "CVF", "FEV1",
		"% de Cambio FEV1 (((FEV1año anterior-FEV1actual)/FEV1año anterior) x100, tener en cuenta que se toman los valores absolutos de los FEV1.", "Interpretación",
		"PA Sistólica (mmHg)", "PA Diastólica (mmHg)", "EKG en reposo", "Prueba de Esfuerzo", "Hipertensión arterial (SI/NO)", "Diabetes (SI/NO)", "Fumador (SI/NO)",
		"Score Framingham", "OD S/C", "OI S/C", "OD C/C", "OI C/C", "OD S/C", "OI S/C", "OD C/C", "OI C/C", "S/C", "C/C", "", "", "OD", "OI", "Peso (K)", "Talla (cm)",
		"IMC", "Dx", "Fototipos", "Hallazgos", "Antecedentes personales", "Antecedentes ocupacionales", "SI", "NO", "", "", "", "", "", "", "", "", "", "", "", "",
		"Minitest psiquiátrico", "Otros (especificar)", "Psicosensómetrico", "Ficha SAHS", "ESTRÉS", "GOLDBERG Ansiedad", "GOLDBERG Depresión", "PERCEPCIÓN DEL RIESGO",
		"Conclusión", "Audit", "Ficha de evaluación neurológica", "Aptitud Altura Estructural", "Test de Impulsividad", "Test Acrofobia", "Test Agorafobia", "", "", "", "", "",
		"MEDICO", "FECHA REVISION", "MEDICO", "FECHA LEV. OBS.", "Cobre en sangre mcg/dL (ug/dL)", "Molibdeno en sangre mcg/dL (ug/dL)", "Plomo inorganico en  sangre mcg/dL (ug/dL)",
		"Cadmio en sangre mcg/L", "", "", "", "Parasitologico Seriado (x 3)", "Coprocultivo", "Hisopado Naso - Faringeo (SEMESTRAL)", "BK esputo", "VDRL", "Fecha INI-FIN",
		"Empresa", "Cargo", "Aptitud", "Cintura", "Cadera", "DESCRIP PUNTJE", "OBSERVACIONES", "Aptitud", "Conducto Auditivo Externo OD", "Membrana Timpánica OD",
		"Conducto Auditivo Externo OI", "Membrana Timpánica OI", "Resultado",
	} {
		header3 = append(header3, cell3)
	}

	if err := streamWriter.SetRow("AB4", header3, excelize.RowOpts{StyleID: styleID}); err != nil {
		fmt.Println(err)
		//return
	}

	merges := map[string]string{
		"A1": "M1", "A2": "A11", "B2": "B11", "C2": "C11", "D2": "D11", "E2": "E11", "F2": "F11", "G2": "G11", "H2": "H11",
		"I2": "I11", "J2": "J11", "K2": "K11", "L2": "L11", "M2": "M11",

		"N1": "AA1", "N2": "N11", "O2": "O11", "P2": "P11", "Q2": "Q11",
		"R2": "R11", "S2": "S11", "T2": "T11", "U2": "U11", "V2": "V11", "W2": "W11", "X2": "X11", "Y2": "Y11", "Z2": "Z11",
		"AA2": "AA11",

		"AB1": "BF1", "AB2": "BF3",
		"AB4": "AB11", "AC4": "AC11", "AD4": "AD11", "AE4": "AE11", "AF4": "AF11", "AG4": "AG11", "AH4": "AH11", "AI4": "AI11", "AJ4": "AJ11", "AK4": "AK11",
		"AL4": "AL11", "AM4": "AM11", "AN4": "AN11", "AO4": "AO11", "AP4": "AP11", "AQ4": "AQ11", "AR4": "AR11", "AS4": "AS11", "AT4": "AT11", "AU4": "AU11",
		"AV4": "AV11", "AW4": "AW11", "AX4": "AX11", "AY4": "AY11", "AZ4": "AZ11",
		"BA4": "BA11", "BB4": "BB11", "BC4": "BC11", "BD4": "BD11", "BE4": "BE11", "BF4": "BF11",

		"BG1": "CY1",
		"BG2": "BG11", "BH2": "BH11", "BI2": "BI11", "BJ2": "BJ11",

		"BK2": "CT2",
		"BK3": "CJ3", "CK3": "CT3",

		"BK4": "BK11", "BL4": "BL11", "BM4": "BM11", "BN4": "BN11", "BO4": "BO11", "BP4": "BP11", "BQ4": "BQ11", "BR4": "BR11", "BS4": "BS11", "BT4": "BT11",
		"BU4": "BU11", "BV4": "BV11", "BW4": "BW11", "BX4": "BX11", "BY4": "BY11", "BZ4": "BZ11", "CA4": "CA11", "CB4": "CB11", "CC4": "CC11", "CD4": "CD11",
		"CE4": "CE11", "CF4": "CF11", "CG4": "CG11", "CH4": "CH11", "CI4": "CI11", "CJ4": "CJ11", "CK4": "CK11", "CL4": "CL11", "CM4": "CM11", "CN4": "CN11",
		"CO4": "CO11", "CP4": "CP11", "CQ4": "CQ11", "CR4": "CR11", "CS4": "CS11", "CT4": "CT11",

		"CU2": "CU11", "CV2": "CY2",

		"CV4": "CV11", "CW4": "CW11", "CX4": "CX11", "CY4": "CY11",
		"CZ2": "CZ11",

		"DA1": "DT1", "DA2": "DJ2", "DK2": "DT2",

		"DA4": "DA11", "DB4": "DB11", "DC4": "DC11", "DD4": "DD11", "DE4": "DE11", "DF4": "DF11", "DG4": "DG11", "DH4": "DH11", "DI4": "DI11", "DJ4": "DJ11",
		"DK4": "DK11", "DL4": "DL11", "DM4": "DM11", "DN4": "DN11", "DO4": "DO11", "DP4": "DP11", "DQ4": "DQ11", "DR4": "DR11", "DS4": "DS11", "DT4": "DT11",

		"DU1": "DW2",

		"DU4": "DU11", "DV4": "DV11", "DW4": "DW11",

		"DX1": "EA2",

		"DX4": "DX11", "DY4": "DY11", "DZ4": "DZ11", "EA4": "EA11",

		"EB1": "EI2",

		"EB4": "EB11", "EC4": "EC11", "ED4": "ED11", "EE4": "EE11", "EF4": "EF11", "EG4": "EG11", "EH4": "EH11", "EI4": "EI11",

		"EJ1": "EW1", "EJ2": "EM2", "EN2": "EQ2", "ER2": "ES2", "ET2": "ET11", "EU2": "EU11", "EV2": "EW2",

		"EJ4": "EJ11", "EK4": "EK11", "EL4": "EL11", "EM4": "EM11", "EN4": "EN11", "EO4": "EO11", "EP4": "EP11", "EQ4": "EQ11", "ER4": "ER11", "ES4": "ES11",
		"EV4": "EV11", "EW4": "EW11",

		"EX1": "FG1", "EX2": "FA2", "FB2": "FG2",
		"EX3": "FA3", "FB3": "FC3", "FD3": "FE3", "FF3": "FG3",

		"EX4": "EX11", "EY4": "EY11", "EZ4": "EZ11", "FA4": "FA11", "FB4": "FB11", "FC4": "FC11", "FD4": "FD11", "FE4": "FE11", "FF4": "FF11", "FG4": "FG11",

		"FH1": "FL1", "FM1": "FQ1",

		"FH2": "FH11", "FI2": "FI11", "FJ2": "FJ11", "FK2": "FK11", "FL2": "FL11", "FM2": "FM11", "FN2": "FN11", "FO2": "FO11", "FP2": "FP11", "FQ2": "FQ11",
		"FR2": "FR11", "FS2": "FS11",

		"FT1": "FU2", "FT3": "FU3",

		"FT4": "FT11", "FU4": "FU11",

		"FV1": "GB2",

		"FV4": "FV11", "FW4": "FW11", "FX4": "FX11", "FY4": "FY11", "FZ4": "FZ11", "GA4": "GA11", "GB4": "GB11",

		"GC1": "GH2",

		"GC4": "GC11", "GD4": "GD11", "GE4": "GE11", "GF4": "GF11", "GG4": "GG11", "GH4": "GH11",

		"GI1": "GI11", "GJ1": "GJ11", "GK1": "GK11", "GL1": "GM1", "GL2": "GL11", "GM2": "GM11",

		"GN1": "GQ1", "GN2": "GO2", "GP2": "GQ2",

		"GN4": "GN11", "GO4": "GO11", "GP4": "GP11", "GQ4": "GQ11",

		"GR1": "GU2", "GR3": "GU3",

		"GR4": "GR11", "GS4": "GS11", "GT4": "GT11", "GU4": "GU11",

		"GV1": "GV11", "GW1": "GW11", "GX1": "GX11",

		"GY1": "HC2",

		"GY4": "GY11", "GZ4": "GZ11", "HA4": "HA11", "HB4": "HB11", "HC4": "HC11",

		"HD1": "HF3",

		"HD4": "HD11", "HE4": "HE11", "HF4": "HF11",

		"HG1": "HG3",

		"HG4": "HG11",

		"HH1": "HI3",

		"HH4": "HH11", "HI4": "HI11",

		"HJ1": "HJ3",

		"HJ4": "HJ11",

		"HK1": "HK3",

		"HK4": "HK11",

		"HL1": "HL3",

		"HL4": "HL11",

		"HM1": "HP3",

		"HM4": "HM11", "HN4": "HN11", "HO4": "HO11", "HP4": "HP11",

		"HQ1": "HQ3",

		"HQ4": "HQ11",
	}

	for x, v := range merges {
		if err := streamWriter.MergeCell(x, v); err != nil {
			fmt.Println(err)
			//return
		}
	}

	//items, _ := c.DB.GetMatrizOnline("2023-05-03", "2023-05-04")
	items, _ := c.DB.GetMatrizOnline(exs.Ini, exs.Fin, exs.OrganizationID)

	y := 12
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
				} else {
					eco = "NO APTO"
				}
			} else {
				eco = "---"
			}
		}

		//---PARTICION NOMBRE----

		partesNombre := strings.Split(e.Name, " ")

		cantidadNombres := len(partesNombre)

		//----------------------------

		//---ELECCION GENERO----

		sex := ""

		if e.SexType == "1" {
			sex = "MASCULINO"
		} else {
			sex = "FEMENINO"
		}

		//----------------------------

		//---PARTIR RESTRICCIONES----

		primerasRestricciones := strings.Split(acuFilter1, ",")

		cantidadRestricciones := len(primerasRestricciones)

		primeraR := "---"
		segundaR := "---"
		terceraR := "---"
		cuartaR := "---"
		quintaR := "---"
		sextaR := "---"

		for i := 0; i < cantidadRestricciones; i++ {
			switch i {
			case 0:
				primeraR = primerasRestricciones[i]
			case 1:
				segundaR = primerasRestricciones[i]
			case 2:
				terceraR = primerasRestricciones[i]
			case 3:
				cuartaR = primerasRestricciones[i]
			case 4:
				quintaR = primerasRestricciones[i]
			case 5:
				sextaR = primerasRestricciones[i]
			}
		}

		//----------------------------

		//---HEMOGRAMA COMPLETO----

		//LEUCOCITOS
		leu, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000424")
		leuV := ""

		for _, leux := range leu {
			leuV = leux.Value1
		}

		//HEMATIES
		hemati, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000422")
		hematiV := ""

		for _, hematix := range hemati {
			hematiV = hematix.Value1
		}

		//HEMOGLOBINA
		hemogl, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001282")
		hemoglV := ""

		for _, hemoglx := range hemogl {
			hemoglV = hemoglx.Value1
		}

		//HEMATOCRITO
		hemato, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		hematoV := ""

		for _, hematox := range hemato {
			hematoV = hematox.Value1
		}

		//NEUTROFILOS - ABASTONADOS %
		neutro, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000426")
		neutroV := ""

		for _, neutrox := range neutro {
			neutroV = neutrox.Value1
		}

		//NEUTROFILOS - ABASTONADOS CEL/ML
		// neutro2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		neutro2V := "NO REALIZAMOS"
		// for _, neutro2x := range neutro2 {
		// 	neutro2V = neutro2x.Value1
		// }

		//PLAQUETAS CEL/ML
		// plaque, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		plaqueV := "NO REALIZAMOS"

		// for _, plaquex := range plaque {
		// 	plaqueV = plaquex.Value1
		// }

		//VCM
		// vcm, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		vcmV := "NO REALIZAMOS"

		// for _, vcmx := range vcm {
		// 	vcmV = vcmx.Value1
		// }

		//HCM
		// hcm, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		hcmV := "NO REALIZAMOS"

		// for _, hcmx := range hcm {
		// 	hcmV = hcmx.Value1
		// }

		//CCMH
		// ccmh, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		ccmhV := "NO REALIZAMOS"

		// for _, ccmhx := range ccmh {
		// 	ccmhV = ccmhx.Value1
		// }

		//RDW
		// rdw, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		rdwV := "NO REALIZAMOS"

		// for _, rdwx := range rdw {
		// 	rdwV = rwdx.Value1
		// }

		//VPM
		// vpm, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		vpmV := "NO REALIZAMOS"

		// for _, vpmx := range vpm {
		// 	vpmV = vpmx.Value1
		// }

		//BASTONES %
		// bastones, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		bastonesV := "NO REALIZAMOS"

		// for _, bastonesx := range bastones {
		// 	bastonesV = bastonesx.Value1
		// }

		//LINFOCITOS %
		linfo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000436")
		linfoV := ""

		for _, linfox := range linfo {
			linfoV = linfox.Value1
		}

		//MONOCITOS %
		mono, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000434")
		monoV := ""

		for _, monox := range mono {
			monoV = monox.Value1
		}

		//EOSINOFILOS %
		eosi, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000430")
		eosiV := ""

		for _, eosix := range eosi {
			eosiV = eosix.Value1
		}

		//BASOFILOS %
		baso, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000432")
		basoV := ""

		for _, basox := range baso {
			basoV = basox.Value1
		}

		//METAMIELOCITOS %
		// metamielocitos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		metamielocitosV := "NO REALIZAMOS"

		// for _, metamielocitosx := range metamielocitos {
		// 	metamielocitosV = metamielocitosx.Value1
		// }

		//MIELOCITOS %
		// mielocitos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		mielocitosV := "NO REALIZAMOS"

		// for _, mielocitosx := range mielocitos {
		// 	mielocitosV = mielocitosx.Value1
		// }

		//PROMIELOCITOS %
		// promie, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		promieV := "NO REALIZAMOS"

		// for _, promiex := range promie {
		// 	promieV = promiex.Value1
		// }

		//BLASTOS %
		// blastos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		blastosV := "NO REALIZAMOS"

		// for _, blastosx := range blastos {
		// 	blastosV = blastosx.Value1
		// }

		//BASTONES CEL/ML
		// bastones2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		bastones2V := "NO REALIZAMOS"

		// for _, bastones2x := range bastones2 {
		// 	bastones2V = bastones2x.Value1
		// }

		//LINFOCITOS CEL/ML
		// linfo2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000436")
		linfo2V := "NO REALIZAMOS"

		// for _, linfo2x := range linfo2 {
		// 	linfo2V = linfo2x.Value1
		// }

		//MONOCITOS CEL/ML
		// mono2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000434")
		mono2V := "NO REALIZAMOS"

		// for _, mono2x := range mono2 {
		// 	mono2V = mono2x.Value1
		// }

		//EOSINOFILOS CEL/ML
		// eosi2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000430")
		eosi2V := "NO REALIZAMOS"

		// for _, eosi2x := range eosi2 {
		// 	eosi2V = eosi2x.Value1
		// }

		//BASOFILOS CEL/ML
		// baso2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000000432")
		baso2V := "NO REALIZAMOS"

		// for _, baso2x := range baso2 {
		// 	baso2V = baso2x.Value1
		// }

		//MIELOCITOS CEL/ML
		// mielocitos2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		mielocitos2V := "NO REALIZAMOS"

		// for _, mielocitos2x := range mielocitos2 {
		// 	mielocitos2V = mielocitos2x.Value1
		// }

		//METAMIELOCITOS CEL/ML
		// metamielocitos2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		metamielocitos2V := "NO REALIZAMOS"

		// for _, metamielocitos2x := range metamielocitos2 {
		// 	metamielocitos2V = metamielocitos2x.Value1
		// }

		//PROMIELOCITOS CEL/ML
		// promie2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		promie2V := "NO REALIZAMOS"

		// for _, promie2x := range promie2 {
		// 	promie2V = promie2x.Value1
		// }

		//BLASTOS CEL/ML
		// blastos2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001280")
		blastos2V := "NO REALIZAMOS"

		// for _, blastos2x := range blastos2 {
		// 	blastos2V = blastos2x.Value1
		// }

		//RENCUENTRO DE PLAQUETAS
		rplaquetas, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000045", "N009-MF000001284")
		rplaquetasV := ""

		for _, rplaquetasx := range rplaquetas {
			rplaquetasV = rplaquetasx.Value1
		}

		//----------------------------------------------------------------

		//---------- GRUPO Y FACTOR SANGUINEO ----------

		//GRUPO
		grupoSan, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000000", "N009-MF000000262")
		grupoSanV := ""

		for _, grupoSanx := range grupoSan {
			grupoSanV = grupoSanx.Value1
		}

		grupoParameter, _ := c.DB.GetValueFromParameterV1("154", grupoSanV)
		grupoDefine := ""

		for _, grupoParameterx := range grupoParameter {
			grupoDefine = grupoParameterx.Value1
		}

		//FACTOR
		factorSan, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000000", "N009-MF000000263")
		factorSanV := ""

		for _, factorSanx := range factorSan {
			factorSanV = factorSanx.Value1
		}

		factorParameter, _ := c.DB.GetValueFromParameterV1("155", factorSanV)
		factorDefine := ""

		for _, factorParameterx := range factorParameter {
			factorDefine = factorParameterx.Value1
		}

		//----------------------------------------------------------------

		//---------- GLUCOSA ----------

		glucosa, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000008", "N009-MF000000261")
		glucosaV := ""

		for _, glucosax := range glucosa {
			glucosaV = glucosax.Value1
		}

		//----------------------------------------------------------------

		//---------- HEMOGLOBINA GLICOSILADA ----------
		// hemoglico, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000008", "N009-MF000000261")
		hemoglicoV := "NO APLICA"

		// for _, hemoglicox := range hemoglico {
		// 	hemoglicoV = hemoglicox.Value1
		// }

		//----------------------------------------------------------------

		//---------- RPR ----------

		rpr, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000003", "N009-MF000000269")
		rprV := ""

		for _, rprx := range rpr {
			rprV = rprx.Value1
		}

		//----------------------------------------------------------------

		//---------- EXAMEN DE ORINA COMPLETO ----------

		//COLOR
		color, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000000444")
		colorV := ""

		for _, colorx := range color {
			colorV = colorx.Value1
		}

		colorParameter, _ := c.DB.GetValueFromParameterV1("257", colorV)
		colorDefine := ""

		for _, colorParameterx := range colorParameter {
			colorDefine = colorParameterx.Value1
		}

		//ASPECTO
		aspecto, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		aspectoV := ""

		for _, aspectox := range aspecto {
			aspectoV = aspectox.Value1
		}

		aspectoParameter, _ := c.DB.GetValueFromParameterV1("258", aspectoV)
		aspectoDefine := ""

		for _, aspectoParameterx := range aspectoParameter {
			aspectoDefine = aspectoParameterx.Value1
		}

		//DENSIDAD
		densidad, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001043")
		densidadV := ""

		for _, densidadx := range densidad {
			densidadV = densidadx.Value1
		}

		//PH
		ph, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001045")
		phV := ""

		for _, phx := range ph {
			phV = phx.Value1
		}

		//GLUCOSA - ORINA
		glucosaOrina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001313")
		glucosaOrinaV := ""

		for _, glucosaOrinax := range glucosaOrina {
			glucosaOrinaV = glucosaOrinax.Value1
		}

		//BILIRRUBINA
		bilirrubina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003439")
		bilirrubinaV := ""

		for _, bilirrubinax := range bilirrubina {
			bilirrubinaV = bilirrubinax.Value1
		}

		//CUERPOS CETONICOS
		cuerposCeto, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001057")
		cuerposCetoV := ""

		for _, cuerposCetox := range cuerposCeto {
			cuerposCetoV = cuerposCetox.Value1
		}

		cuerposCetoParameter, _ := c.DB.GetValueFromParameterV1("265", cuerposCetoV)
		cuerposCetoDefine := ""

		for _, cuerposCetoParameterx := range cuerposCetoParameter {
			cuerposCetoDefine = cuerposCetoParameterx.Value1
		}

		//PROTEINAS
		proteinas, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001053")
		proteinasV := ""

		for _, proteinasx := range proteinas {
			proteinasV = proteinasx.Value1
		}

		proteinasParameter, _ := c.DB.GetValueFromParameterV1("261", proteinasV)
		proteinasDefine := ""

		for _, proteinasParameterx := range proteinasParameter {
			proteinasDefine = proteinasParameterx.Value1
		}

		//UROBILINOGENO
		urobili, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001049")
		urobiliV := ""

		for _, urobilix := range urobili {
			urobiliV = urobilix.Value1
		}

		urobiliParameter, _ := c.DB.GetValueFromParameterV1("264", urobiliV)
		urobiliDefine := ""

		for _, urobiliParameterx := range urobiliParameter {
			urobiliDefine = urobiliParameterx.Value1
		}

		//NITRITOS
		nitrito, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001055")
		nitritoV := ""

		for _, nitritox := range nitrito {
			nitritoV = nitritox.Value1
		}

		nitritoParameter, _ := c.DB.GetValueFromParameterV1("260", nitritoV)
		nitritoDefine := ""

		for _, nitritoParameterx := range nitritoParameter {
			nitritoDefine = nitritoParameterx.Value1
		}

		//HEMOGLOBINA - ORINA
		// hemoOrina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		hemoOrinaV := "NO REALIZAMOS"

		// for _, hemoOrinax := range hemoOrina {
		// 	hemoOrinaV = hemoOrinax.Value1
		// }

		//LEUCOCITOS - ORINA
		leucoOrina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003438")
		leucoOrinaV := ""

		for _, leucoOrinax := range leucoOrina {
			leucoOrinaV = leucoOrinax.Value1
		}

		//HEMATIES - ORINA
		hemaOrina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001063")
		hemaOrinaV := ""

		for _, hemaOrinax := range hemaOrina {
			hemaOrinaV = hemaOrinax.Value1
		}

		//CELULAS EPITELIALES
		celulasEpi, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		celulasEpiV := ""

		for _, celulasEpix := range celulasEpi {
			celulasEpiV = celulasEpix.Value1
		}

		//LEVADURAS
		levadura, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003437")
		levaduraV := ""

		for _, levadurax := range levadura {
			levaduraV = levadurax.Value1
		}

		levaduraParameter, _ := c.DB.GetValueFromParameterV1("269", levaduraV)
		levaduraDefine := ""

		for _, levaduraParameterx := range levaduraParameter {
			levaduraDefine = levaduraParameterx.Value1
		}

		//CRISTALES
		// cristales, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		cristalesV := "NO REALIZAMOS"

		// for _, cristalesx := range cristales {
		// 	cristalesV = cristalesx.Value1
		// }

		// cristalesParameter, _ := c.DB.GetValueFromParameterV1("258", cristalesV)
		// cristalesDefine := ""

		// for _, cristalesParameterx := range cristalesParameter {
		// 	cristalesDefine = cristalesParameterx.Value1
		// }

		//C. ACURICO
		// CAcuri, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		CAcuriV := "NO REALIZAMOS"

		// for _, CAcurix := range CAcuri {
		// 	CAcuriV = CAcurix.Value1
		// }

		// CAcuriParameter, _ := c.DB.GetValueFromParameterV1("258", CAcuriV)
		// CAcuriDefine := ""

		// for _, CAcuriParameterx := range CAcuriParameter {
		// 	CAcuriDefine = CAcuriParameterx.Value1
		// }

		//C. FOSFAMORFOS
		CFostamorfos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003243")
		CFostamorfosV := ""

		for _, CFostamorfosx := range CFostamorfos {
			CFostamorfosV = CFostamorfosx.Value1
		}

		CFostamorfosParameter, _ := c.DB.GetValueFromParameterV1("270", CFostamorfosV)
		CFostamorfosDefine := ""

		for _, CFostamorfosParameterx := range CFostamorfosParameter {
			CFostamorfosDefine = CFostamorfosParameterx.Value1
		}

		//C. URATAMORFOS
		// CUratamorfos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		CUratamorfosV := "NO REALIZAMOS"

		// for _, CUratamorfosx := range CUratamorfos {
		// 	CUratamorfosV = CUratamorfosx.Value1
		// }

		// CUratamorfosParameter, _ := c.DB.GetValueFromParameterV1("258", CUratamorfosV)
		// CUratamorfosDefine := ""

		// for _, CUratamorfosParameterx := range CUratamorfosParameter {
		// 	CUratamorfosDefine = CUratamorfosParameterx.Value1
		// }

		//C. OXCALCIO
		COxcalcio, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001065")
		COxcalcioV := ""

		for _, COxcalciox := range COxcalcio {
			COxcalcioV = COxcalciox.Value1
		}

		COxcalcioParameter, _ := c.DB.GetValueFromParameterV1("270", COxcalcioV)
		COxcalcioDefine := ""

		for _, COxcalcioParameterx := range COxcalcioParameter {
			COxcalcioDefine = COxcalcioParameterx.Value1
		}

		//C. FOSF-TRIPLES
		CFostTriples, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003244")
		CFostTriplesV := ""

		for _, CFostTriplesx := range CFostTriples {
			CFostTriplesV = CFostTriplesx.Value1
		}

		CFostTriplesParameter, _ := c.DB.GetValueFromParameterV1("270", CFostTriplesV)
		CFostTriplesDefine := ""

		for _, CFostTriplesParameterx := range CFostTriplesParameter {
			CFostTriplesDefine = CFostTriplesParameterx.Value1
		}

		//CILINDROS
		// cilindros, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		cilindrosV := "NO REALIZAMOS"

		// for _, cilindrosx := range cilindros {
		// 	cilindrosV = cilindrosx.Value1
		// }

		// cilindrosParameter, _ := c.DB.GetValueFromParameterV1("258", cilindrosV)
		// cilindrosDefine := ""

		// for _, cilindrosParameterx := range cilindrosParameter {
		// 	cilindrosDefine = cilindrosParameterx.Value1
		// }

		//CILINDROS HIALINOS
		CHialinos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001069")
		CHialinosV := ""

		for _, CHialinosx := range CHialinos {
			CHialinosV = CHialinosx.Value1
		}

		//GRANULOSOS
		granulosos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		granulososV := ""

		for _, granulososx := range granulosos {
			granulososV = granulososx.Value1
		}

		//FILAMENTOS MUCOIDEOS
		FilamentoMucoideo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000003440")
		FilamentoMucoideoV := ""

		for _, FilamentoMucoideox := range FilamentoMucoideo {
			FilamentoMucoideoV = FilamentoMucoideox.Value1
		}

		FilamentoMucoideoParameter, _ := c.DB.GetValueFromParameterV1("269", FilamentoMucoideoV)
		FilamentoMucoideoDefine := ""

		for _, FilamentoMucoideoParameterx := range FilamentoMucoideoParameter {
			FilamentoMucoideoDefine = FilamentoMucoideoParameterx.Value1
		}

		//GERMENES
		germenes, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001067")
		germenesV := ""

		for _, germenesx := range germenes {
			germenesV = germenesx.Value1
		}

		germenesParameter, _ := c.DB.GetValueFromParameterV1("269", germenesV)
		germenesDefine := ""

		for _, germenesParameterx := range germenesParameter {
			germenesDefine = germenesParameterx.Value1
		}

		//----------------------------------------------------------------

		//---------- TOXICOLOGICO ----------

		//COCAINA

		cocaina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000633", "N009-MF000005173")
		cocainaV := ""

		for _, cocainax := range cocaina {
			cocainaV = cocainax.Value1
		}

		cocainaParameter, _ := c.DB.GetValueFromParameterV1("305", cocainaV)
		cocainaDefine := ""

		for _, cocainaParameterx := range cocainaParameter {
			cocainaDefine = cocainaParameterx.Value1
		}

		toxcocaina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000461", "N009-MF000003740")
		toxcocainaV := ""

		for _, toxcocainax := range toxcocaina {
			toxcocainaV = toxcocainax.Value1
		}

		toxcocainaParameter, _ := c.DB.GetValueFromParameterV1("305", toxcocainaV)
		toxcocainaDefine := ""

		for _, toxcocainaParameterx := range toxcocainaParameter {
			toxcocainaDefine = toxcocainaParameterx.Value1
		}

		CocainaDefinitiva := ""

		if cocainaDefine == "" {
			CocainaDefinitiva = toxcocainaDefine
		}
		if toxcocainaDefine == "" {
			CocainaDefinitiva = cocainaDefine
		}

		//MARIHUANA
		marihuana, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000634", "N009-MF000005174")
		marihuanaV := ""

		for _, marihuanax := range marihuana {
			marihuanaV = marihuanax.Value1
		}

		marihuanaParameter, _ := c.DB.GetValueFromParameterV1("305", marihuanaV)
		marihuanaDefine := ""

		for _, marihuanaParameterx := range marihuanaParameter {
			marihuanaDefine = marihuanaParameterx.Value1
		}

		toxmarihuana, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000461", "N009-MF000003739")
		toxmarihuanaV := ""

		for _, toxmarihuanax := range toxmarihuana {
			toxmarihuanaV = toxmarihuanax.Value1
		}

		toxmarihuanaParameter, _ := c.DB.GetValueFromParameterV1("305", toxmarihuanaV)
		toxmarihuanaDefine := ""

		for _, toxmarihuanaParameterx := range toxmarihuanaParameter {
			toxmarihuanaDefine = toxmarihuanaParameterx.Value1
		}

		MarihuanaDefinitiva := ""

		if marihuanaDefine == "" {
			MarihuanaDefinitiva = toxmarihuanaDefine
		}
		if toxmarihuanaDefine == "" {
			MarihuanaDefinitiva = marihuanaDefine
		}

		//ANFETAMINA
		anfetamina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000043", "N009-MF000000391")
		anfetaminaV := ""

		for _, anfetaminax := range anfetamina {
			anfetaminaV = anfetaminax.Value1
		}

		anfetaminaParameter, _ := c.DB.GetValueFromParameterV1("305", anfetaminaV)
		anfetaminaDefine := ""

		for _, anfetaminaParameterx := range anfetaminaParameter {
			anfetaminaDefine = anfetaminaParameterx.Value1
		}

		//METANFETAMINA
		metanfetamina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000419", "N009-MF000003215")
		metanfetaminaV := ""

		for _, metanfetaminax := range metanfetamina {
			metanfetaminaV = metanfetaminax.Value1
		}

		metanfetaminaParameter, _ := c.DB.GetValueFromParameterV1("305", metanfetaminaV)
		metanfetaminaDefine := ""

		for _, metanfetaminaParameterx := range metanfetaminaParameter {
			metanfetaminaDefine = metanfetaminaParameterx.Value1
		}

		//METADONA
		metadona, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000418", "N009-MF000003214")
		metadonaV := ""

		for _, metadonax := range metadona {
			metadonaV = metadonax.Value1
		}

		metadonaParameter, _ := c.DB.GetValueFromParameterV1("305", metadonaV)
		metadonaDefine := ""

		for _, metadonaParameterx := range metadonaParameter {
			metadonaDefine = metadonaParameterx.Value1
		}

		//MORFINA
		morfina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000420", "N009-MF000003216")
		morfinaV := ""

		for _, morfinax := range morfina {
			morfinaV = morfinax.Value1
		}

		morfinaParameter, _ := c.DB.GetValueFromParameterV1("305", morfinaV)
		morfinaDefine := ""

		for _, morfinaParameterx := range morfinaParameter {
			morfinaDefine = morfinaParameterx.Value1
		}

		//FENCICLIDINA
		fenciclidina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000423", "N009-MF000003219")
		fenciclidinaV := ""

		for _, fenciclidinax := range fenciclidina {
			fenciclidinaV = fenciclidinax.Value1
		}

		fenciclidinaParameter, _ := c.DB.GetValueFromParameterV1("305", fenciclidinaV)
		fenciclidinaDefine := ""

		for _, fenciclidinaParameterx := range fenciclidinaParameter {
			fenciclidinaDefine = fenciclidinaParameterx.Value1
		}

		//BARBITURICOS
		barbituricos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000417", "N009-MF000003213")
		barbituricosV := ""

		for _, barbituricosx := range barbituricos {
			barbituricosV = barbituricosx.Value1
		}

		barbituricosParameter, _ := c.DB.GetValueFromParameterV1("305", barbituricosV)
		barbituricosDefine := ""

		for _, barbituricosParameterx := range barbituricosParameter {
			barbituricosDefine = barbituricosParameterx.Value1
		}

		//BENZODIACEPINAS
		benzodiacepina, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000040", "N009-MF000000395")
		benzodiacepinaV := ""

		for _, benzodiacepinax := range benzodiacepina {
			benzodiacepinaV = benzodiacepinax.Value1
		}

		benzodiacepinaParameter, _ := c.DB.GetValueFromParameterV1("305", benzodiacepinaV)
		benzodiacepinaDefine := ""

		for _, benzodiacepinaParameterx := range benzodiacepinaParameter {
			benzodiacepinaDefine = benzodiacepinaParameterx.Value1
		}

		if benzodiacepinaDefine == "" {
			benzodiacepinaDefine = "NO APLICA"
		}

		//ANTIDEPRESIVOS
		// antidepresivos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000046", "N009-MF000001041")
		antidepresivosV := "NO APLICA"

		// for _, antidepresivosx := range antidepresivos {
		// 	antidepresivosV = antidepresivosx.Value1
		// }

		// antidepresivosParameter, _ := c.DB.GetValueFromParameterV1("258", antidepresivosV)
		// antidepresivosDefine := ""

		// for _, antidepresivosParameterx := range antidepresivosParameter {
		// 	antidepresivosDefine = antidepresivosParameterx.Value1
		// }

		//----------------------------------------------------------------

		//---------- Prueba de HCGB (Edad Fertil  - Mujeres) mIU/mL ----------

		pruebaEmbarazo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000027", "N009-MF000000270")
		pruebaEmbarazoV := ""

		for _, pruebaEmbarazox := range pruebaEmbarazo {
			pruebaEmbarazoV = pruebaEmbarazox.Value1
		}

		pruebaEmbarazoParameter, _ := c.DB.GetValueFromParameterV1("203", pruebaEmbarazoV)
		pruebaEmbarazoDefine := ""

		for _, pruebaEmbarazoParameterx := range pruebaEmbarazoParameter {
			pruebaEmbarazoDefine = pruebaEmbarazoParameterx.Value1
		}

		//----------------------------------------------------------------

		//---------- Perfil lipídico ----------

		//COLESTEROL
		colesterol, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000114", "N009-MF000001904")
		colesterolV := ""

		for _, colesterolx := range colesterol {
			colesterolV = colesterolx.Value1
		}

		colesterolTotal, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000016", "N009-MF000001086")
		colesterolTotalV := ""

		for _, colesterolTotalx := range colesterolTotal {
			colesterolTotalV = colesterolTotalx.Value1
		}

		colesterolDefinitivo := ""

		if colesterolV == "" {
			colesterolDefinitivo = colesterolTotalV
		}
		if colesterolTotalV == "" {
			colesterolDefinitivo = colesterolV
		}
		if colesterolTotalV == "" && colesterolV == "" {
			colesterolDefinitivo = "NO APLICA"
		}

		//HDL
		Hdl, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000114", "N009-MF000000254")
		HdlV := ""

		for _, Hdlx := range Hdl {
			HdlV = Hdlx.Value1
		}

		if HdlV == "" {
			HdlV = "NO APLICA"
		}

		//LDL
		ldl, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000114", "N009-MF000001073")
		ldlV := ""

		for _, ldlx := range ldl {
			ldlV = ldlx.Value1
		}

		if ldlV == "" {
			ldlV = "NO APLICA"
		}

		//TRIGLICERIDOS
		trigliceridos, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000114", "N009-MF000001906")
		trigliceridosV := ""

		for _, trigliceridosx := range trigliceridos {
			trigliceridosV = trigliceridosx.Value1
		}

		trigliceridos2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000017", "N009-MF000001296")
		trigliceridos2V := ""

		for _, trigliceridos2x := range trigliceridos2 {
			trigliceridos2V = trigliceridos2x.Value1
		}

		trigliceridosDefinitivos := ""

		if trigliceridosV == "" {
			trigliceridosDefinitivos = trigliceridos2V
		}
		if trigliceridos2V == "" {
			trigliceridosDefinitivos = trigliceridosV
		}
		if trigliceridosV == "" && trigliceridos2V == "" {
			trigliceridosDefinitivos = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- ODONTOLOGIA ----------

		NroCaries, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000027", "N002-ODO00000193")
		NroCariesV := ""

		for _, NroCariesx := range NroCaries {
			NroCariesV = NroCariesx.Value1
		}

		// Dividir la cadena en un slice de strings
		NroCariesNum := strings.Split(NroCariesV, ";")

		// Contar la cantidad de elementos en el slice
		NroCariesDefinitivo := len(NroCariesNum)

		//----------------------------------------------------------------

		//---------- AUDIOMETRIA ----------

		//OD
		OD500, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000001")
		OD500V := ""

		for _, OD500x := range OD500 {
			OD500V = OD500x.Value1
		}

		OD1000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000002")
		OD1000V := ""

		for _, OD1000x := range OD1000 {
			OD1000V = OD1000x.Value1
		}

		OD2000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000003")
		OD2000V := ""

		for _, OD2000x := range OD2000 {
			OD2000V = OD2000x.Value1
		}

		OD3000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000004")
		OD3000V := ""

		for _, OD3000x := range OD3000 {
			OD3000V = OD3000x.Value1
		}

		OD4000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000005")
		OD4000V := ""

		for _, OD4000x := range OD4000 {
			OD4000V = OD4000x.Value1
		}

		OD6000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000006")
		OD6000V := ""

		for _, OD6000x := range OD6000 {
			OD6000V = OD6000x.Value1
		}

		OD8000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000007")
		OD8000V := ""

		for _, OD8000x := range OD8000 {
			OD8000V = OD8000x.Value1
		}

		STSODV := "NO REALIZAMOS"

		InterClinicaODV := "- - -"

		InterOcupacionalOD, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000005", "N009-MF000004461")
		InterOcupacionalODV := ""

		for _, InterOcupacionalODx := range InterOcupacionalOD {
			InterOcupacionalODV = InterOcupacionalODx.Value1
		}

		InterOcupacionalODParameter, _ := c.DB.GetValueFromParameterV1("301", InterOcupacionalODV)
		InterOcupacionalODDefine := ""

		for _, InterOcupacionalODParameterx := range InterOcupacionalODParameter {
			InterOcupacionalODDefine = InterOcupacionalODParameterx.Value1
		}

		//OI
		OI500, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000015")
		OI500V := ""

		for _, OI500x := range OI500 {
			OI500V = OI500x.Value1
		}

		OI1000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000016")
		OI1000V := ""

		for _, OI1000x := range OI1000 {
			OI1000V = OI1000x.Value1
		}

		OI2000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000017")
		OI2000V := ""

		for _, OI2000x := range OI2000 {
			OI2000V = OI2000x.Value1
		}

		OI3000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000018")
		OI3000V := ""

		for _, OI3000x := range OI3000 {
			OI3000V = OI3000x.Value1
		}

		OI4000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000019")
		OI4000V := ""

		for _, OI4000x := range OI4000 {
			OI4000V = OI4000x.Value1
		}

		OI6000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000020")
		OI6000V := ""

		for _, OI6000x := range OI6000 {
			OI6000V = OI6000x.Value1
		}

		OI8000, _ := c.DB.GetValueCustomerV2(e.VServiceid, "N002-ME000000005", "N002-AUD00000021")
		OI8000V := ""

		for _, OI8000x := range OI8000 {
			OI8000V = OI8000x.Value1
		}

		STSOIV := "NO REALIZAMOS"

		InterClinicaOIV := "- - -"

		InterOcupacionalOI, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000005", "N009-MF000004462")
		InterOcupacionalOIV := ""

		for _, InterOcupacionalOIx := range InterOcupacionalOI {
			InterOcupacionalOIV = InterOcupacionalOIx.Value1
		}

		InterOcupacionalOIParameter, _ := c.DB.GetValueFromParameterV1("302", InterOcupacionalOIV)
		InterOcupacionalOIDefine := ""

		for _, InterOcupacionalOIParameterx := range InterOcupacionalOIParameter {
			InterOcupacionalOIDefine = InterOcupacionalOIParameterx.Value1
		}

		//BILATERAL
		Bilateral, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000005", "N009-MF000004463")
		BilateralV := ""

		for _, Bilateralx := range Bilateral {
			BilateralV = Bilateralx.Value1
		}

		BilateralParameter, _ := c.DB.GetValueFromParameterV1("303", BilateralV)
		BilateralDefine := ""

		for _, BilateralParameterx := range BilateralParameter {
			BilateralDefine = BilateralParameterx.Value1
		}

		if InterOcupacionalODDefine == "- - -" {
			InterOcupacionalODDefine = BilateralDefine
		}
		if InterOcupacionalOIDefine == "- - -" {
			InterOcupacionalOIDefine = BilateralDefine
		}

		//----------------------------------------------------------------

		//---------- RAYOS X ----------

		//OIT
		ToraxOIT, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000062", "N009-MF000002142")
		ToraxOITV := ""

		for _, ToraxOITx := range ToraxOIT {
			ToraxOITV = ToraxOITx.Value1
		}

		ToraxOITParameter, _ := c.DB.GetValueFromParameterV1("318", ToraxOITV)
		ToraxOITDefine := ""

		for _, ToraxOITParameterx := range ToraxOITParameter {
			ToraxOITDefine = ToraxOITParameterx.Value1
		}

		if ToraxOITDefine == "" {
			ToraxOITDefine = "NO APLICA"
		}

		//HALLAZGOS - OIT
		HallazgosOIT, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000062", "N009-MF000001039")
		HallazgosOITV := ""

		for _, HallazgosOITx := range HallazgosOIT {
			HallazgosOITV = HallazgosOITx.Value1
		}

		//HALLAZGOS - TORAX
		HallazgosTORAXN, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000032", "N009-MF000002134")
		HallazgosTORAXNV := ""

		for _, HallazgosTORAXNx := range HallazgosTORAXN {
			HallazgosTORAXNV = HallazgosTORAXNx.Value1
		}

		HallazgosTORAXA, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000032", "N009-MF000004528")
		HallazgosTORAXAV := ""

		for _, HallazgosTORAXAx := range HallazgosTORAXA {
			HallazgosTORAXAV = HallazgosTORAXAx.Value1
		}

		HallazgoToraxDefinitivo := ""

		if HallazgosTORAXNV == "1" {
			HallazgoToraxDefinitivo = "PLACA NORMAL"
		}
		if HallazgosTORAXAV == "1" {
			HallazgoToraxDefinitivo = "PLACA ANORMAL"
		}

		//HALLAZGO - TOTAL
		HallazgoToraxTotal := ""

		if HallazgosOITV == "" {
			HallazgoToraxTotal = HallazgoToraxDefinitivo
		}
		if HallazgoToraxDefinitivo == "" {
			HallazgoToraxTotal = HallazgosOITV
		}

		//DX-TORAX
		dxTorax, _ := c.DB.GetDxSingle(e.VServiceid, "N002-ME000000032")
		dxToraxV := ""

		for _, dxToraxx := range dxTorax {
			dxToraxV = dxToraxx.Name
		}

		//DX-OIT
		dxOIT, _ := c.DB.GetDxSingle(e.VServiceid, "N009-ME000000062")
		dxOITV := ""

		for _, dxOITx := range dxOIT {
			dxOITV = dxOITx.Name
		}

		//DX - RAYOS - DEFINITIVO
		DxRayosDefinitivo := ""

		if dxToraxV == "" {
			DxRayosDefinitivo = dxOITV
		}

		if dxOITV == "" {
			DxRayosDefinitivo = dxToraxV
		}

		//----------------------------------------------------------------

		//---------- ESPIROMETRIA ----------

		//CVF
		CVF, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000031", "N002-MF000000286")
		CVFV := ""

		for _, CVFx := range CVF {
			CVFV = CVFx.Value1
		}

		if CVFV == "" {
			CVFV = "NO APLICA"
		}

		//FEV1
		FEV1, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000031", "N002-MF000000287")
		FEV1V := ""

		for _, FEV1x := range FEV1 {
			FEV1V = FEV1x.Value1
		}

		if FEV1V == "" {
			FEV1V = "NO APLICA"
		}

		//% de Cambio FEV1
		PercentFEV1V := "- - -"

		//DX - ESPIRO
		DxEspiro, _ := c.DB.GetDxSingle(e.VServiceid, "N002-ME000000031")
		DxEspiroV := ""

		for _, DxEspirox := range DxEspiro {
			DxEspiroV = DxEspirox.Name
		}

		if DxEspiroV == "" {
			DxEspiroV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- CARDIOLOGIA ----------

		//PA SISTOLICA
		PASis, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000001", "N002-MF000000001")
		PASisV := ""

		for _, PASisx := range PASis {
			PASisV = PASisx.Value1
		}

		//PA DIASTOLICA
		PADias, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000001", "N002-MF000000002")
		PADiasV := ""

		for _, PADiasx := range PADias {
			PADiasV = PADiasx.Value1
		}

		//EKG REPOSO
		EKGReposo, _ := c.DB.GetDxSingle(e.VServiceid, "N002-ME000000025")
		EKGReposoV := ""

		for _, EKGReposox := range EKGReposo {
			EKGReposoV = EKGReposox.Name
		}

		if EKGReposoV == "" {
			EKGReposoV = "NO APLICA"
		}

		//PRUEBA DE ESFUERZO
		PruebaEsfuerzo, _ := c.DB.GetDxSingle(e.VServiceid, "N002-ME000000029")
		PruebaEsfuerzoV := ""

		for _, PruebaEsfuerzox := range PruebaEsfuerzo {
			PruebaEsfuerzoV = PruebaEsfuerzox.Name
		}

		if PruebaEsfuerzoV == "1" {
			PruebaEsfuerzoV = "PRUEBA DE ESFUERZO SIN ALTERACION"
		}
		if PruebaEsfuerzoV == "" {
			PruebaEsfuerzoV = "NO APLICA"
		}

		//HIPERTENSION
		HiperTension, _ := c.DB.GetCheckDx(e.VServiceid, "N009-DD000000436")
		HiperTensionV := ""

		for _, HiperTensionx := range HiperTension {
			HiperTensionV = HiperTensionx.Name
		}

		if HiperTensionV != "" {
			HiperTensionV = "SI"
		}
		if HiperTensionV == "" {
			HiperTensionV = "NO"
		}

		//DIABETES
		Diabetes, _ := c.DB.GetCheckDx(e.VServiceid, "N009-DD000000642")
		DiabetesV := ""

		for _, Diabetesx := range Diabetes {
			DiabetesV = Diabetesx.Name
		}

		if DiabetesV != "" {
			DiabetesV = "SI"
		}
		if DiabetesV == "" {
			DiabetesV = "NO"
		}

		//FUMADOR
		Fumador, _ := c.DB.GetNoxiusHabitats(e.VServiceid, "1")
		FumadorV := ""

		for _, Fumadorx := range Fumador {
			FumadorV = Fumadorx.Frequency
		}

		if FumadorV != "NADA" {
			FumadorV = "SI"
		}
		if FumadorV == "NADA" {
			FumadorV = "NO"
		}

		//SCORE FRAMINGHAM
		ScoreFraminghamV := "NO REALIZAMOS"

		//----------------------------------------------------------------

		//---------- OFTALMOLOGIA ----------

		//CERCA ODSC - SIMPLE
		CercaODSCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003547")
		CercaODSCSV := ""

		for _, CercaODSCSx := range CercaODSCS {
			CercaODSCSV = CercaODSCSx.Value1
		}

		CercaODSCSCercaODSCSrameter, _ := c.DB.GetValueFromParameterV1("290", CercaODSCSV)
		CercaODSCSDefine := ""

		for _, CercaODSCSCercaODSCSrameterx := range CercaODSCSCercaODSCSrameter {
			CercaODSCSDefine = CercaODSCSCercaODSCSrameterx.Value1
		}

		//CERCA ODSC - COMPLETO
		CercaODSCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003569")
		CercaODSCCV := ""

		for _, CercaODSCCx := range CercaODSCC {
			CercaODSCCV = CercaODSCCx.Value1
		}

		CercaODSCCCercaODSCCrameter, _ := c.DB.GetValueFromParameterV1("290", CercaODSCCV)
		CercaODSCCDefine := ""

		for _, CercaODSCCCercaODSCCrameterx := range CercaODSCCCercaODSCCrameter {
			CercaODSCCDefine = CercaODSCCCercaODSCCrameterx.Value1
		}

		//CERCA ODSC - DEFINITIVO
		CercaODSC := ""

		if CercaODSCSDefine == "" {
			CercaODSC = CercaODSCCDefine
		} else {
			CercaODSC = CercaODSCSDefine
		}

		//CERCA OISC - SIMPLE
		CercaOISCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003548")
		CercaOISCSV := ""

		for _, CercaOISCSx := range CercaOISCS {
			CercaOISCSV = CercaOISCSx.Value1
		}

		CercaOISCSCercaOISCSrameter, _ := c.DB.GetValueFromParameterV1("290", CercaOISCSV)
		CercaOISCSDefine := ""

		for _, CercaOISCSCercaOISCSrameterx := range CercaOISCSCercaOISCSrameter {
			CercaOISCSDefine = CercaOISCSCercaOISCSrameterx.Value1
		}

		//CERCA OISC - COMPLETO
		CercaOISCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003570")
		CercaOISCCV := ""

		for _, CercaOISCCx := range CercaOISCC {
			CercaOISCCV = CercaOISCCx.Value1
		}

		CercaOISCCCercaOISCCrameter, _ := c.DB.GetValueFromParameterV1("290", CercaOISCCV)
		CercaOISCCDefine := ""

		for _, CercaOISCCCercaOISCCrameterx := range CercaOISCCCercaOISCCrameter {
			CercaOISCCDefine = CercaOISCCCercaOISCCrameterx.Value1
		}

		//CERCA OISC - DEFINITIVO
		CercaOISC := ""

		if CercaOISCSDefine == "" {
			CercaOISC = CercaOISCCDefine
		} else {
			CercaOISC = CercaOISCSDefine
		}

		//CERCA ODCC - SIMPLE
		CercaODCCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003549")
		CercaODCCSV := ""

		for _, CercaODCCSx := range CercaODCCS {
			CercaODCCSV = CercaODCCSx.Value1
		}

		CercaODCCSCercaODCCSrameter, _ := c.DB.GetValueFromParameterV1("290", CercaODCCSV)
		CercaODCCSDefine := ""

		for _, CercaODCCSCercaODCCSrameterx := range CercaODCCSCercaODCCSrameter {
			CercaODCCSDefine = CercaODCCSCercaODCCSrameterx.Value1
		}

		//CERCA ODCC - COMPLETO
		CercaODCCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003571")
		CercaODCCCV := ""

		for _, CercaODCCCx := range CercaODCCC {
			CercaODCCCV = CercaODCCCx.Value1
		}

		CercaODCCCCercaODCCCrameter, _ := c.DB.GetValueFromParameterV1("290", CercaODCCCV)
		CercaODCCCDefine := ""

		for _, CercaODCCCCercaODCCCrameterx := range CercaODCCCCercaODCCCrameter {
			CercaODCCCDefine = CercaODCCCCercaODCCCrameterx.Value1
		}

		//CERCA ODCC - DEFINITIVO
		CercaODCC := ""

		if CercaODCCSDefine == "" {
			CercaODCC = CercaODCCCDefine
		} else {
			CercaODCC = CercaODCCSDefine
		}

		//CERCA OICC - SIMPLE
		CercaOICCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003550")
		CercaOICCSV := ""

		for _, CercaOICCSx := range CercaOICCS {
			CercaOICCSV = CercaOICCSx.Value1
		}

		CercaOICCSCercaOICCSrameter, _ := c.DB.GetValueFromParameterV1("290", CercaOICCSV)
		CercaOICCSDefine := ""

		for _, CercaOICCSCercaOICCSrameterx := range CercaOICCSCercaOICCSrameter {
			CercaOICCSDefine = CercaOICCSCercaOICCSrameterx.Value1
		}

		//CERCA OICC - COMPLETO
		CercaOICCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003572")
		CercaOICCCV := ""

		for _, CercaOICCCx := range CercaOICCC {
			CercaOICCCV = CercaOICCCx.Value1
		}

		CercaOICCCCercaOICCCrameter, _ := c.DB.GetValueFromParameterV1("290", CercaOICCCV)
		CercaOICCCDefine := ""

		for _, CercaOICCCCercaOICCCrameterx := range CercaOICCCCercaOICCCrameter {
			CercaOICCCDefine = CercaOICCCCercaOICCCrameterx.Value1
		}

		//CERCA OICC - DEFINITIVO
		CercaOICC := ""

		if CercaOICCSDefine == "" {
			CercaOICC = CercaOICCCDefine
		} else {
			CercaOICC = CercaOICCSDefine
		}

		//LEJOS ODSC - SIMPLE
		LejosODSCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003551")
		LejosODSCSV := ""

		for _, LejosODSCSx := range LejosODSCS {
			LejosODSCSV = LejosODSCSx.Value1
		}

		LejosODSCSLejosODSCSrameter, _ := c.DB.GetValueFromParameterV1("287", LejosODSCSV)
		LejosODSCSDefine := ""

		for _, LejosODSCSLejosODSCSrameterx := range LejosODSCSLejosODSCSrameter {
			LejosODSCSDefine = LejosODSCSLejosODSCSrameterx.Value1
		}

		//LEJOS ODSC - COMPLETO
		LejosODSCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003565")
		LejosODSCCV := ""

		for _, LejosODSCCx := range LejosODSCC {
			LejosODSCCV = LejosODSCCx.Value1
		}

		LejosODSCCLejosODSCCrameter, _ := c.DB.GetValueFromParameterV1("287", LejosODSCCV)
		LejosODSCCDefine := ""

		for _, LejosODSCCLejosODSCCrameterx := range LejosODSCCLejosODSCCrameter {
			LejosODSCCDefine = LejosODSCCLejosODSCCrameterx.Value1
		}

		//LEJOS ODSC - DEFINITIVO
		LejosODSC := ""

		if LejosODSCSDefine == "" {
			LejosODSC = LejosODSCCDefine
		} else {
			LejosODSC = LejosODSCSDefine
		}

		//LEJOS OISC - SIMPLE
		LejosOISCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003552")
		LejosOISCSV := ""

		for _, LejosOISCSx := range LejosOISCS {
			LejosOISCSV = LejosOISCSx.Value1
		}

		LejosOISCSLejosOISCSrameter, _ := c.DB.GetValueFromParameterV1("287", LejosOISCSV)
		LejosOISCSDefine := ""

		for _, LejosOISCSLejosOISCSrameterx := range LejosOISCSLejosOISCSrameter {
			LejosOISCSDefine = LejosOISCSLejosOISCSrameterx.Value1
		}

		//LEJOS OISC - COMPLETO
		LejosOISCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003566")
		LejosOISCCV := ""

		for _, LejosOISCCx := range LejosOISCC {
			LejosOISCCV = LejosOISCCx.Value1
		}

		LejosOISCCLejosOISCCrameter, _ := c.DB.GetValueFromParameterV1("287", LejosOISCCV)
		LejosOISCCDefine := ""

		for _, LejosOISCCLejosOISCCrameterx := range LejosOISCCLejosOISCCrameter {
			LejosOISCCDefine = LejosOISCCLejosOISCCrameterx.Value1
		}

		//LEJOS OISC - DEFINITIVO
		LejosOISC := ""

		if LejosOISCSDefine == "" {
			LejosOISC = LejosOISCCDefine
		} else {
			LejosOISC = LejosOISCSDefine
		}

		//LEJOS ODCC - SIMPLE
		LejosODCCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003553")
		LejosODCCSV := ""

		for _, LejosODCCSx := range LejosODCCS {
			LejosODCCSV = LejosODCCSx.Value1
		}

		LejosODCCSLejosODCCSrameter, _ := c.DB.GetValueFromParameterV1("287", LejosODCCSV)
		LejosODCCSDefine := ""

		for _, LejosODCCSLejosODCCSrameterx := range LejosODCCSLejosODCCSrameter {
			LejosODCCSDefine = LejosODCCSLejosODCCSrameterx.Value1
		}

		//LEJOS ODCC - COMPLETO
		LejosODCCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003567")
		LejosODCCCV := ""

		for _, LejosODCCCx := range LejosODCCC {
			LejosODCCCV = LejosODCCCx.Value1
		}

		LejosODCCCLejosODCCCrameter, _ := c.DB.GetValueFromParameterV1("287", LejosODCCCV)
		LejosODCCCDefine := ""

		for _, LejosODCCCLejosODCCCrameterx := range LejosODCCCLejosODCCCrameter {
			LejosODCCCDefine = LejosODCCCLejosODCCCrameterx.Value1
		}

		//LEJOS ODCC - DEFINITIVO
		LejosODCC := ""

		if LejosODCCSDefine == "" {
			LejosODCC = LejosODCCCDefine
		} else {
			LejosODCC = LejosODCCSDefine
		}

		//LEJOS OICC - SIMPLE
		LejosOICCS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003554")
		LejosOICCSV := ""

		for _, LejosOICCSx := range LejosOICCS {
			LejosOICCSV = LejosOICCSx.Value1
		}

		LejosOICCSLejosOICCSrameter, _ := c.DB.GetValueFromParameterV1("287", LejosOICCSV)
		LejosOICCSDefine := ""

		for _, LejosOICCSLejosOICCSrameterx := range LejosOICCSLejosOICCSrameter {
			LejosOICCSDefine = LejosOICCSLejosOICCSrameterx.Value1
		}

		//LEJOS OICC - COMPLETO
		LejosOICCC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003568")
		LejosOICCCV := ""

		for _, LejosOICCCx := range LejosOICCC {
			LejosOICCCV = LejosOICCCx.Value1
		}

		LejosOICCCLejosOICCCrameter, _ := c.DB.GetValueFromParameterV1("287", LejosOICCCV)
		LejosOICCCDefine := ""

		for _, LejosOICCCLejosOICCCrameterx := range LejosOICCCLejosOICCCrameter {
			LejosOICCCDefine = LejosOICCCLejosOICCCrameterx.Value1
		}

		//LEJOS OICC - DEFINITIVO
		LejosOICC := ""

		if LejosOICCSDefine == "" {
			LejosOICC = LejosOICCCDefine
		} else {
			LejosOICC = LejosOICCSDefine
		}

		//AGUDEZA BINOCULAR
		ABinocularSCV := "NO APLICA"

		ABinocularCCV := "NO APLICA"

		//ISHIHARA SIMPLE
		ISHIHARAS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003555")
		ISHIHARASV := ""

		for _, ISHIHARASx := range ISHIHARAS {
			ISHIHARASV = ISHIHARASx.Value1
		}

		ISHIHARASISHIHARASrameter, _ := c.DB.GetValueFromParameterV1("217", ISHIHARASV)
		ISHIHARASDefine := ""

		for _, ISHIHARASISHIHARASrameterx := range ISHIHARASISHIHARASrameter {
			ISHIHARASDefine = ISHIHARASISHIHARASrameterx.Value1
		}

		//ISHIHARA COMPLETO
		ISHIHARAC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000003573")
		ISHIHARACV := ""

		for _, ISHIHARACx := range ISHIHARAC {
			ISHIHARACV = ISHIHARACx.Value1
		}

		ISHIHARACISHIHARACrameter, _ := c.DB.GetValueFromParameterV1("217", ISHIHARACV)
		ISHIHARACDefine := ""

		for _, ISHIHARACISHIHARACrameterx := range ISHIHARACISHIHARACrameter {
			ISHIHARACDefine = ISHIHARACISHIHARACrameterx.Value1
		}

		//ISHIHARA DEFINE
		ISHIHARADefinitivo := ""

		if ISHIHARASDefine == "" {
			ISHIHARADefinitivo = ISHIHARACDefine
		} else {
			ISHIHARADefinitivo = ISHIHARASDefine
		}

		//ESTEROSCOPICA SIMPLE
		EsteroscopiaS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000447", "N009-MF000004316")
		EsteroscopiaSV := ""

		for _, EsteroscopiaSx := range EsteroscopiaS {
			EsteroscopiaSV = EsteroscopiaSx.Value1
		}

		EsteroscopiaSEsteroscopiaSrameter, _ := c.DB.GetValueFromParameterV1("221", EsteroscopiaSV)
		EsteroscopiaSDefine := ""

		for _, EsteroscopiaSEsteroscopiaSrameterx := range EsteroscopiaSEsteroscopiaSrameter {
			EsteroscopiaSDefine = EsteroscopiaSEsteroscopiaSrameterx.Value1
		}

		//ESTEROSCOPICA COMPLETO
		EsteroscopiaC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000448", "N009-MF000003576")
		EsteroscopiaCV := ""

		for _, EsteroscopiaCx := range EsteroscopiaC {
			EsteroscopiaCV = EsteroscopiaCx.Value1
		}

		//ESTEROSCOPICA DEFINITVO
		EsteroscopicaDefinitivo := ""

		if EsteroscopiaSDefine == "" {
			EsteroscopicaDefinitivo = EsteroscopiaCV
		} else {
			EsteroscopicaDefinitivo = EsteroscopiaSDefine
		}

		//CAMPOS VISUALES
		CamposVisualesSCV := "NO APLICA"

		CamposVisualesCCV := "NO APLICA"

		//----------------------------------------------------------------

		//---------- NUTRICION ----------

		//PESO
		Peso, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000002", "N002-MF000000008")
		PesoV := ""

		for _, Pesox := range Peso {
			PesoV = Pesox.Value1
		}

		//TALLA
		Talla, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000002", "N002-MF000000008")
		TallaV := ""

		for _, Tallax := range Talla {
			TallaV = Tallax.Value1
		}

		//IMC
		IMC, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000002", "N002-MF000000009")
		IMCV := ""

		for _, IMCx := range IMC {
			IMCV = IMCx.Value1
		}

		//DX - NUTRICIONAL
		DxNutricional, _ := c.DB.GetDxSingle(e.VServiceid, "N002-ME000000002")
		DxNutricionalV := ""

		for _, DxNutricionalx := range DxNutricional {
			DxNutricionalV = DxNutricionalx.Name
		}

		//----------------------------------------------------------------

		//---------- ANTECEDENTES ----------

		//FOTOTIPO
		FasciogramaV := "NO APLICA"

		//HALLAZGOS
		HallazgosAnteV := "NO APLICA"

		//ANTECEDENTES PERSONALES
		DxPersonal, _ := c.DB.GetAntecedentesPersonales(e.VPersonId)
		DxPersonalV := ""

		for _, DxPersonalx := range DxPersonal {
			if DxPersonalV == "" {
				DxPersonalV = DxPersonalx.DxDetail
			} else {
				DxPersonalV = DxPersonalV + ", " + DxPersonalx.DxDetail
			}
		}

		if DxPersonalV == "" {
			DxPersonalV = "NINGUNO"
		}

		//ANTECEDENTES OCUPACIONALES
		AnteOcupacionalV := "NINGUNO"

		//ALERGIAS
		Alergias, _ := c.DB.GetCheckDx(e.VServiceid, "N009-DD000000633")

		AlergiasSI := ""
		AlergiasNO := ""

		for _, Alergiasx := range Alergias {
			if Alergiasx.Name == "" {
				AlergiasNO = "NO"
			} else {
				AlergiasSI = "SI"
			}

		}

		//----------------------------------------------------------------

		//---------- INMUNIZACIONES ----------

		//TETANOS
		Tetanos1, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000604", "N009-MF000004672")
		Tetanos1V := ""

		for _, Tetanos1x := range Tetanos1 {
			Tetanos1V = "1ERA: " + Tetanos1x.Value1
		}

		Tetanos2, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000604", "N009-MF000004674")
		Tetanos2V := ""

		for _, Tetanos2x := range Tetanos2 {
			Tetanos2V = "2DA :" + Tetanos2x.Value1
		}

		Tetanos3, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000604", "N009-MF000004676")
		Tetanos3V := ""

		for _, Tetanos3x := range Tetanos3 {
			Tetanos3V = "3ERA :" + Tetanos3x.Value1
		}

		Tetanos4, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000604", "N009-MF000004678")
		Tetanos4V := ""

		for _, Tetanos4x := range Tetanos4 {
			Tetanos4V = "4TA :" + Tetanos4x.Value1
		}

		Tetanos5, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000604", "N009-MF000004680")
		Tetanos5V := ""

		for _, Tetanos5x := range Tetanos5 {
			Tetanos5V = "5TA :" + Tetanos5x.Value1
		}

		TetanosDefinitivo := ""

		if Tetanos1V == "" {
			TetanosDefinitivo = "NO APLICA"
		} else {
			TetanosDefinitivo = Tetanos1V + ", " + Tetanos2V + ", " + Tetanos3V + ", " + Tetanos4V + ", " + Tetanos5V
		}

		//FIEBRE TIFOIDEA
		FiebreTifoideaV := "NO APLICA"

		//HEPATITIS A
		Hepatitis1A, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000659", "N009-MF000005439")
		Hepatitis1AV := ""

		for _, Hepatitis1Ax := range Hepatitis1A {
			Hepatitis1AV = "1ERA :" + Hepatitis1Ax.Value1
		}

		Hepatitis2A, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000659", "N009-MF000005441")
		Hepatitis2AV := ""

		for _, Hepatitis2Ax := range Hepatitis2A {
			Hepatitis2AV = "2DA :" + Hepatitis2Ax.Value1
		}

		Hepatitis3A, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000659", "N009-MF000005443")
		Hepatitis3AV := ""

		for _, Hepatitis3Ax := range Hepatitis3A {
			Hepatitis3AV = "3ERA :" + Hepatitis3Ax.Value1
		}

		HepatitisADefinitivo := ""

		if Hepatitis1AV == "" {
			HepatitisADefinitivo = "NO APLICA"
		} else {
			HepatitisADefinitivo = Hepatitis1AV + ", " + Hepatitis2AV + ", " + Hepatitis3AV
		}

		//HEPATITIS B
		Hepatitis1B, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000613", "N009-MF000004641")
		Hepatitis1BV := ""

		for _, Hepatitis1Bx := range Hepatitis1B {
			Hepatitis1BV = "1ERA :" + Hepatitis1Bx.Value1
		}

		Hepatitis2B, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000613", "N009-MF000004642")
		Hepatitis2BV := ""

		for _, Hepatitis2Bx := range Hepatitis2B {
			Hepatitis2BV = "2DA :" + Hepatitis2Bx.Value1
		}

		Hepatitis3B, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000613", "N009-MF000004643")
		Hepatitis3BV := ""

		for _, Hepatitis3Bx := range Hepatitis3B {
			Hepatitis3BV = "3ERA :" + Hepatitis3Bx.Value1
		}

		HepatitisBDefinitivo := ""

		if Hepatitis1BV == "" {
			HepatitisBDefinitivo = "NO APLICA"
		} else {
			HepatitisBDefinitivo = Hepatitis1BV + ", " + Hepatitis2BV + ", " + Hepatitis3BV
		}

		//INFLUENZA
		Influenza, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000603", "N009-MF000004639")
		InfluenzaV := ""

		for _, Influenzax := range Influenza {
			InfluenzaV = Influenzax.Value1
		}

		if InfluenzaV == "" {
			InfluenzaV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- COVID ----------
		COVID, _ := c.DB.GetCheckAntePer(e.VPersonId, "N009-DD000001698")
		COVIDV := ""

		for _, COVIDx := range COVID {
			COVIDV = COVIDx.DxDetail
		}

		// Separar por ":"
		partes := strings.Split(COVIDV, "\n")

		// Extraer las fechas y almacenarlas en variables
		var primeraDosis, segundaDosis, terceraDosis, cuartaDosis, quintaDosis string

		for _, parte := range partes {
			if strings.Contains(parte, "/") {
				partesDosisFecha := strings.Split(parte, ":")
				if len(partesDosisFecha) == 2 { // Ensure there are two parts after split
					dosis := strings.TrimSpace(partesDosisFecha[0])
					fecha := strings.TrimSpace(partesDosisFecha[1])

					switch dosis {
					case "1° DOSIS", "1ERA DOSIS":
						primeraDosis = fecha
					case "2° DOSIS", "2DA DOSIS":
						segundaDosis = fecha
					case "3° DOSIS", "3ERA DOSIS":
						terceraDosis = fecha
					case "4° DOSIS", "4TA DOSIS":
						cuartaDosis = fecha
					case "5° DOSIS", "5TA DOSIS":
						quintaDosis = fecha
					}
				} else {
					fmt.Println("Warning: Invalid format in line:", parte)
				}
			}
		}
		//----------------------------------------------------------------

		//---------- DESCARTE TBC ----------

		//CUADRO CLINICO
		CuadroClinicoSR, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000022", "N009-RES00000009")
		CuadroClinicoSRV := ""

		for _, CuadroClinicoSRx := range CuadroClinicoSR {
			CuadroClinicoSRV = CuadroClinicoSRx.Value1
		}

		if CuadroClinicoSRV == "" {
			CuadroClinicoSRV = "NO APLICA"
		}

		//BACILOSCOPIA
		Baciloscopia, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000022", "N009-RES00000008")
		BaciloscopiaV := ""

		for _, Baciloscopiax := range Baciloscopia {
			BaciloscopiaV = Baciloscopiax.Value1
		}

		if BaciloscopiaV == "" {
			BaciloscopiaV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- PSICOLOGIA ----------

		//MINI TEST PSIQUIATRICO
		MiniTestPsiquiatricoV := ""

		//OTROS TEST
		OtrosTest, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000033", "N002-MF000000278")
		OtrosTestV := ""

		for _, OtrosTestx := range OtrosTest {
			OtrosTestV = OtrosTestx.Value1
		}

		if OtrosTestV == "" {
			OtrosTestV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- PERFIL CONDUCTOR ----------

		//PSICOSENSOMETRICO
		Psicosensometrico, _ := c.DB.GetDxSingle(e.VServiceid, "N009-ME000000617")
		PsicosensometricoV := ""

		for _, Psicosensometricox := range Psicosensometrico {
			PsicosensometricoV = Psicosensometricox.Name
		}

		//FICHA SAHS
		SAS, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000435", "N009-ME000003311")
		SASV := ""

		for _, SASx := range SAS {
			SASV = SASx.Value1
		}

		SASParameter, _ := c.DB.GetValueFromParameterV1("111", SASV)
		SASDefine := ""

		for _, SASParameterx := range SASParameter {
			SASDefine = SASParameterx.Value1
		}

		SASDefinitivo := ""

		if SASDefine == "SI" {
			SASDefinitivo = "NORMAL"
		} else if SASDefine == "NO" {
			SASDefinitivo = "ANORMAL"
		} else if SASDefine == "" {
			SASDefinitivo = "NO APLICA"
		}

		//ESTRES
		Estres, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N100-ME000000470", "N009-ME000003311")
		EstresV := ""

		for _, Estresx := range Estres {
			EstresV = Estresx.Value1
		}

		EstresDefinitivo := ""

		if EstresV == "" {
			EstresDefinitivo = "NO APLICA"
		} else {
			EstresDefinitivo = "NIVEL DE ESTRES: " + EstresV
		}

		//GOLDBERG - ANSIEDAD
		GoldbergAnsiedad, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000621", "N009-MF000004881")
		GoldbergAnsiedadV := ""

		for _, GoldbergAnsiedadx := range GoldbergAnsiedad {
			GoldbergAnsiedadV = GoldbergAnsiedadx.Value1
		}

		GoldbergAnsiedadDefinitivo := ""

		if GoldbergAnsiedadV == "" {
			GoldbergAnsiedadDefinitivo = "NO APLICA"
		} else {
			GoldbergAnsiedadDefinitivo = "ANSIEDAD: " + GoldbergAnsiedadV
		}

		//GOLDBERG - DEPRESION
		GoldbergDepresion, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000621", "N009-MF000004881")
		GoldbergDepresionV := ""

		for _, GoldbergDepresionx := range GoldbergDepresion {
			GoldbergDepresionV = GoldbergDepresionx.Value1
		}

		GoldbergDepresionDefinitivo := ""

		if GoldbergDepresionV == "" {
			GoldbergDepresionDefinitivo = "NO APLICA"
		} else {
			GoldbergDepresionDefinitivo = "ANSIEDAD: " + GoldbergDepresionV
		}

		//PERCEPCION DE RIESGO
		PercepcionRiesgoV := ""

		//CONCLUSION
		ConclusionConductorV := "NINGUNA"

		//----------------------------------------------------------------

		//---------- TRABAJOS EN ALTURA ESTRUCTURAL 1.8 ----------

		//AUDIT
		Audit, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000532", "N009-MF000005078")
		AuditV := ""

		for _, Auditx := range Audit {
			AuditV = Auditx.Value1
		}

		AuditDefinitivo := ""

		if AuditV == "" {
			AuditDefinitivo = "NO APLICA"
		} else {
			AuditDefinitivo = "TEST DE AUDIT: " + AuditV
		}

		//EVALUACION NEUROLOGICA
		Neurologico, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000022", "N009-MF000000619")
		NeurologicoV := ""

		for _, Neurologicox := range Neurologico {
			NeurologicoV = Neurologicox.Value1
		}

		//APTITUD ALTURA ESTRUCTURAL
		AlturaEstructural, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000015", "N009-MF000000357")
		AlturaEstructuralV := ""

		for _, AlturaEstructuralx := range AlturaEstructural {
			AlturaEstructuralV = AlturaEstructuralx.Value1
		}

		//TEST DE IMPULSIVIDAD
		TestImpulsividadV := ""

		//TEST DE ACROFOBIA
		Acrofobia, _ := c.DB.GetDxSingle(e.VServiceid, "N009-ME000000618")
		AcrofobiaV := ""
		i := 0

		for _, Acrofobiax := range Acrofobia {
			if i == 0 {
				AcrofobiaV = Acrofobiax.Name
			} else {
				AcrofobiaV = Acrofobiax.Name + ", " + Acrofobiax.Name
			}
			i++
		}

		if AcrofobiaV == "" {
			AcrofobiaV = "NO APLICA"
		}

		//TEST DE AGORAFOBIA
		TestAgorafobiaV := "NO APLICA"

		//----------------------------------------------------------------

		//---------- EVALUACION NEUROLOGICA POR MEDICO NEUROLOGO ----------
		EvaNeurologicaNeurologoV := "NO APLICA"

		//----------------------------------------------------------------

		//---------- EEG ----------
		ElectroEncefaloGramaV := "NO APLICA"

		//----------------------------------------------------------------

		//---------- MANIPULADOR DE ALIMENTOS Y RESIDUOS ----------
		ManipuladorAlimentosV := "- - -"

		//----------------------------------------------------------------

		//---------- DATOS DEL CMA ----------
		ClinicaV := "CENTRO MEDICO HOLOSALUD"

		DoctorAcargoV := "DORA ESTELA CANALES GUILLEN"

		//----------------------------------------------------------------

		//---------- DATOS DEL UMC ----------
		FRMedicoV := "- - -"

		FRMOV := "- - -"

		LevantamientoMedicoV := "- - -"

		LevantamientoFRV := "- - -"

		//----------------------------------------------------------------

		//---------- METALES PESADOS ----------

		//COBRE
		Cobre, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000426", "N009-MF000003222")
		CobreV := ""

		for _, Cobrex := range Cobre {
			CobreV = Cobrex.Value1
		}

		if CobreV == "" {
			CobreV = "NA"
		}

		//MOLIBDENO
		MolibdenoV := "NA"

		//PLOMO
		Plomo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000060", "N009-MF000001158")
		PlomoV := ""

		for _, Plomox := range Plomo {
			PlomoV = Plomox.Value1
		}

		if PlomoV == "" {
			PlomoV = "NA"
		}

		//CADMIO
		Cadmio, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N010-ME000000395", "N009-MF000005370")
		CadmioV := ""

		for _, Cadmiox := range Cadmio {
			CadmioV = Cadmiox.Value1
		}

		if CadmioV == "" {
			CadmioV = "NA"
		}

		//----------------------------------------------------------------

		//---------- SATISFACCION DEL USUARIO ----------

		SatisfaccionUsuarioV := ""

		ProyectoTrabajadorV := e.Location

		ClinicaOrigenV := "CENTRO MEDICO HOLOSALUD"

		//----------------------------------------------------------------

		//---------- MANIPULADOR DE ALIMENTOS ----------

		//PARASITOLOGICO
		Parasitologico, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000049", "N009-MF000003261")
		ParasitologicoV := ""

		for _, Parasitologicox := range Parasitologico {
			ParasitologicoV = Parasitologicox.Value1
		}

		if ParasitologicoV == "" {
			ParasitologicoV = "NO APLICA"
		}

		//COPROCULTIVO
		Coprocultivo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000658", "N009-MF000005431")
		CoprocultivoV := ""

		for _, Coprocultivox := range Coprocultivo {
			CoprocultivoV = Coprocultivox.Value1
		}

		if CoprocultivoV == "" {
			CoprocultivoV = "NO APLICA"
		}

		//HISOPADO NASOFARINGEO
		HisopadoNasofaringeo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000118", "N009-MF000002030")
		HisopadoNasofaringeoV := ""

		for _, HisopadoNasofaringeox := range HisopadoNasofaringeo {
			HisopadoNasofaringeoV = HisopadoNasofaringeox.Value1
		}

		if HisopadoNasofaringeoV == "" {
			HisopadoNasofaringeoV = "NO APLICA"
		}

		//BK ESPUTO
		BKEsputo, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000081", "N009-MF000001373")
		BKEsputoV := ""

		for _, BKEsputox := range BKEsputo {
			BKEsputoV = BKEsputox.Value1
		}

		if BKEsputoV == "" {
			BKEsputoV = "NO APLICA"
		}

		//VDRL
		VDRL, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000003", "N009-MF000000269")
		VDRLV := ""

		for _, VDRLx := range VDRL {
			VDRLV = VDRLx.Value1
		}

		if VDRLV == "" {
			VDRLV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- HISTORIA OCUPACIONAL ----------

		//FECHA INICIO - FECHA FIN
		FechaIniFechaFinV := e.ServiceDate + " - " + e.ExpirationDate

		//EMPRESA
		EmpresaOrgV := e.OrgName

		//CARGO
		CargoV := e.PersonOcupation

		//----------------------------------------------------------------

		//---------- EVA. ESPACIOS CONFINADOS ----------
		EspaciosConfinadosApto, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000436", "N009-MF000003359")
		EspaciosConfinadosAptoV := ""

		for _, EspaciosConfinadosAptox := range EspaciosConfinadosApto {
			EspaciosConfinadosAptoV = EspaciosConfinadosAptox.Value1
		}

		EspaciosConfinadosNoApto, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N009-ME000000436", "N009-MF000003360")
		EspaciosConfinadosNoAptoV := ""

		for _, EspaciosConfinadosNoAptox := range EspaciosConfinadosNoApto {
			EspaciosConfinadosNoAptoV = EspaciosConfinadosNoAptox.Value1
		}

		EspaciosAptitudDefnitivo := ""

		if EspaciosConfinadosAptoV == "" || EspaciosConfinadosNoAptoV == "" {
			EspaciosAptitudDefnitivo = "NO APLICA"
		} else if EspaciosConfinadosAptoV != "0" {
			EspaciosAptitudDefnitivo = "APTO"
		} else if EspaciosConfinadosNoAptoV != "0" {
			EspaciosAptitudDefnitivo = "NO APTO"
		}
		//----------------------------------------------------------------

		//---------- EVA. ANTROPOMETRICA ----------

		//CINTURA
		IndiceCintura, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000002", "N002-MF000000012")
		IndiceCinturaV := ""

		for _, IndiceCinturax := range IndiceCintura {
			IndiceCinturaV = IndiceCinturax.Value1
		}

		if IndiceCinturaV == "" {
			IndiceCinturaV = "NO APLICA"
		}

		//CADERA
		PerCadera, _ := c.DB.GetValueCustomerV1(e.VServiceid, "N002-ME000000002", "N002-MF000000011")
		PerCaderaV := ""

		for _, PerCaderax := range PerCadera {
			PerCaderaV = PerCaderax.Value1
		}

		if PerCaderaV == "" {
			PerCaderaV = "NO APLICA"
		

		//----------------------------------------------------------------

		//---------- TEST EPWORTH ----------
		TestEpworth, _ := c.DB.GetDxSingle(e.VServiceid, "N009-ME000000216")
		TestEpworthV := ""
		i := 0

		for _, TestEpworthx := range TestEpworth {
			TestEpworthV = TestEpworthx.Name
		}

		if TestEpworthV == "" {
			TestEpworthV = "NO APLICA"
		}

		//----------------------------------------------------------------

		//---------- RPR ----------

		//----------------------------------------------------------------

		//---------- RPR ----------

		//----------------------------------------------------------------

		//RowCellValue = append(make([]interface{}, 0), strconv.Itoa(x), e.DocNumber, e.Ape2, e.Ape1, e.Name, e.Name, strconv.Itoa(Age(ages)), e.EsoName, e.ProtocolName, e.ServiceDate[0:10], e.PersonOcupation, e.Aptitude, acuFilter1, "NO APLICA", "1", arr1[0], arr2[0], "1", arr1[1], arr2[1], "1", arr1[2], arr2[2], alti, eco, "EVALUADO")

		if cantidadNombres == 1 {
			RowCellValue = append(make([]interface{}, 0), strconv.Itoa(x), e.DocNumber, e.Ape2, e.Ape1, partesNombre[0], "", e.Bithdate,
				strconv.Itoa(Age(ages)), sex, e.Birthplace, e.Direccion, e.OrgName, e.PersonOcupation, e.EsoName, e.ServiceDate[0:10], e.Aptitude, e.ExpirationDate[0:10],
				"", "", "", "", primeraR, segundaR, terceraR, cuartaR, quintaR, sextaR, leuV, hematiV, hemoglV, hematoV, neutroV, neutro2V, plaqueV, vcmV, hcmV,
				ccmhV, rdwV, vpmV, bastonesV, linfoV, monoV, eosiV, basoV, metamielocitosV, mielocitosV, promieV, blastosV, bastones2V, linfo2V, mono2V, eosi2V,
				baso2V, mielocitos2V, metamielocitos2V, promie2V, blastos2V, rplaquetasV, grupoDefine+" - "+factorDefine, glucosaV, hemoglicoV, rprV,
				colorDefine, aspectoDefine, densidadV, phV, glucosaOrinaV, bilirrubinaV, cuerposCetoDefine, proteinasDefine, urobiliDefine, nitritoDefine,
				hemoOrinaV, leucoOrinaV, hemaOrinaV, celulasEpiV, levaduraDefine, cristalesV, CAcuriV, CFostamorfosDefine, CUratamorfosV, COxcalcioDefine,
				CFostTriplesDefine, cilindrosV, CHialinosV, granulososV, FilamentoMucoideoDefine, germenesDefine, CocainaDefinitiva, MarihuanaDefinitiva,
				anfetaminaDefine, metanfetaminaDefine, metadonaDefine, morfinaDefine, fenciclidinaDefine, barbituricosDefine, benzodiacepinaDefine, antidepresivosV,
				pruebaEmbarazoDefine, colesterolDefinitivo, HdlV, ldlV, trigliceridosDefinitivos, NroCariesDefinitivo, OD500V, OD1000V, OD2000V, OD3000V, OD4000V,
				OD6000V, OD8000V, STSODV, InterClinicaODV, InterOcupacionalODDefine, OI500V, OI1000V, OI2000V, OI3000V, OI4000V, OI6000V, OI8000V, STSOIV,
				InterClinicaOIV, InterOcupacionalOIDefine, ToraxOITDefine, HallazgoToraxTotal, DxRayosDefinitivo, CVFV, FEV1V, PercentFEV1V, DxEspiroV,
				PASisV, PADiasV, EKGReposoV, PruebaEsfuerzoV, HiperTensionV, DiabetesV, FumadorV, ScoreFraminghamV, CercaODSC, CercaOISC, CercaODCC, CercaOICC,
				LejosODSC, LejosOISC, LejosODCC, LejosOICC, ABinocularSCV, ABinocularCCV, ISHIHARADefinitivo, CamposVisualesSCV, CamposVisualesCCV,
				EsteroscopicaDefinitivo, PesoV, TallaV, IMCV, DxNutricionalV, FasciogramaV, HallazgosAnteV, DxPersonalV, AnteOcupacionalV, AlergiasSI,
				AlergiasNO, TetanosDefinitivo, FiebreTifoideaV, HepatitisADefinitivo, HepatitisBDefinitivo, InfluenzaV, primeraDosis, segundaDosis, terceraDosis,
				cuartaDosis, quintaDosis, CuadroClinicoSRV, BaciloscopiaV, MiniTestPsiquiatricoV, OtrosTestV, PsicosensometricoV, SASDefinitivo, EstresDefinitivo,
				GoldbergAnsiedadDefinitivo, GoldbergDepresionDefinitivo, PercepcionRiesgoV, ConclusionConductorV, AuditDefinitivo, NeurologicoV, AlturaEstructuralV,
				TestImpulsividadV, AcrofobiaV, TestAgorafobiaV, EvaNeurologicaNeurologoV, ElectroEncefaloGramaV, ManipuladorAlimentosV, ClinicaV, DoctorAcargoV,
				FRMedicoV, FRMOV, LevantamientoMedicoV, LevantamientoFRV, CobreV, MolibdenoV, PlomoV, CadmioV, SatisfaccionUsuarioV, ProyectoTrabajadorV,
				ClinicaOrigenV, ParasitologicoV, CoprocultivoV, HisopadoNasofaringeoV, BKEsputoV, VDRLV, FechaIniFechaFinV, EmpresaOrgV, CargoV, 
				EspaciosAptitudDefnitivo, IndiceCinturaV, PerCaderaV, TestEpworthV, eco, alti, acuFilter1)

		} else if cantidadNombres > 1 {

			RowCellValue = append(make([]interface{}, 0), strconv.Itoa(x), e.DocNumber, e.Ape2, e.Ape1, partesNombre[0], partesNombre[1], e.Bithdate,
				strconv.Itoa(Age(ages)), sex, e.Birthplace, e.Direccion, e.OrgName, e.PersonOcupation, e.EsoName, e.ServiceDate[0:10], e.Aptitude, e.ExpirationDate[0:10],
				"", "", "", "", primeraR, segundaR, terceraR, cuartaR, quintaR, sextaR, leuV, hematiV, hemogl, hematoV, neutroV, neutro2V, plaqueV, vcmV, hcmV,
				ccmhV, rdwV, vpmV, bastonesV, linfoV, monoV, eosiV, basoV, metamielocitosV, mielocitosV, promieV, blastosV, bastones2V, linfo2V, mono2V, eosi2V,
				baso2V, mielocitos2V, metamielocitos2V, promie2V, blastos2V, rplaquetasV, grupoDefine+" - "+factorDefine, glucosaV, hemoglicoV, rprV,
				colorDefine, aspectoDefine, densidadV, phV, glucosaOrinaV, bilirrubinaV, cuerposCetoDefine, proteinasDefine, urobiliDefine, nitritoDefine,
				hemoOrinaV, leucoOrinaV, hemaOrinaV, celulasEpiV, levaduraDefine, cristalesV, CAcuriV, CFostamorfosDefine, CUratamorfosV, COxcalcioDefine,
				CFostTriplesDefine, cilindrosV, CHialinosV, granulososV, FilamentoMucoideoDefine, germenesDefine, CocainaDefinitiva, MarihuanaDefinitiva,
				anfetaminaDefine, metanfetaminaDefine, metadonaDefine, morfinaDefine, fenciclidinaDefine, barbituricosDefine, benzodiacepinaDefine, antidepresivosV,
				pruebaEmbarazoDefine, colesterolDefinitivo, HdlV, ldlV, trigliceridosDefinitivos, NroCariesDefinitivo, OD500V, OD1000V, OD2000V, OD3000V, OD4000V,
				OD6000V, OD8000V, STSODV, InterClinicaODV, InterOcupacionalODDefine, OI500V, OI1000V, OI2000V, OI3000V, OI4000V, OI6000V, OI8000V, STSOIV,
				InterClinicaOIV, InterOcupacionalOIDefine, ToraxOITDefine, HallazgoToraxTotal, DxRayosDefinitivo, CVFV, FEV1V, PercentFEV1V, DxEspiroV,
				PASisV, PADiasV, EKGReposoV, PruebaEsfuerzoV, HiperTensionV, DiabetesV, FumadorV, ScoreFraminghamV, CercaODSC, CercaOISC, CercaODCC, CercaOICC,
				LejosODSC, LejosOISC, LejosODCC, LejosOICC, ABinocularSCV, ABinocularCCV, ISHIHARADefinitivo, CamposVisualesSCV, CamposVisualesCCV,
				EsteroscopicaDefinitivo, PesoV, TallaV, IMCV, DxNutricionalV, FasciogramaV, HallazgosAnteV, DxPersonalV, AnteOcupacionalV, AlergiasSI,
				AlergiasNO, TetanosDefinitivo, FiebreTifoideaV, HepatitisADefinitivo, HepatitisBDefinitivo, InfluenzaV, primeraDosis, segundaDosis, terceraDosis,
				cuartaDosis, quintaDosis, CuadroClinicoSRV, BaciloscopiaV, MiniTestPsiquiatricoV, OtrosTestV, PsicosensometricoV, SASDefinitivo, EstresDefinitivo,
				GoldbergAnsiedadDefinitivo, GoldbergDepresionDefinitivo, PercepcionRiesgoV, ConclusionConductorV, AuditDefinitivo, NeurologicoV, AlturaEstructuralV,
				TestImpulsividadV, AcrofobiaV, TestAgorafobiaV, EvaNeurologicaNeurologoV, ElectroEncefaloGramaV, ManipuladorAlimentosV, ClinicaV, DoctorAcargoV,
				FRMedicoV, FRMOV, LevantamientoMedicoV, LevantamientoFRV, CobreV, MolibdenoV, PlomoV, CadmioV, SatisfaccionUsuarioV, ProyectoTrabajadorV,
				ClinicaOrigenV, ParasitologicoV, CoprocultivoV, HisopadoNasofaringeoV, BKEsputoV, VDRLV, FechaIniFechaFinV, EmpresaOrgV, CargoV, 
				EspaciosAptitudDefnitivo, IndiceCinturaV, PerCaderaV, TestEpworthV, eco, alti, acuFilter1)
		}

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

	/*
		if err := f.SetRowHeight("Sheet1", 2, 45); err != nil {
			fmt.Println(err)
			//return
		}

		if err := f.SetRowHeight("Sheet1", 4, 178.5); err != nil {
			fmt.Println(err)
			//return
		}
	*/
	if err := f.SaveAs("\\\\HOLO-SERVIDOR\\archivos sistema_2\\TEMPORAL\\" + exs.OrganizationID + ".xlsx"); err != nil {
		println(err.Error())
	}

	var filePath string

	filePath = "\\\\HOLO-SERVIDOR\\archivos sistema_2\\TEMPORAL\\" + exs.OrganizationID + ".xlsx"

	if len(filePath) == 0 {
		return "", errors.New("no aceptado")
	}

	return filePath, nil
}
