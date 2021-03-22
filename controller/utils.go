package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func checkError(err error, ctx string) {
	if err != nil {
		log.Println("Error en: " + ctx)
		log.Println(err)
	}
}

func returnErr(w http.ResponseWriter, err error, operation string) {
	_, _ = fmt.Fprintln(w, fmt.Sprintf("Hubo un error al %s, error: %s", operation, err.Error()))
}

func splitDate(r rune) bool {
	return r == '-' || r == 'T'
}

func formatDate(date string) string {
	data := strings.FieldsFunc(date, splitDate)
	return data[2] + data[1] + data[0]
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
