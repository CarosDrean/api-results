package controller

import "log"

func checkError(err error, ctx string) {
	if err != nil {
		log.Println("Error en: " + ctx)
		log.Println(err)
	}
}
