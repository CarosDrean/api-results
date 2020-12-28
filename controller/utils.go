package controller

import (
	"fmt"
	"log"
	"net/http"
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
