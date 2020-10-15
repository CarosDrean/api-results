package controller

import (
	"net/http"
	"path"
)

func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("\\\\DESKTOP-QD7QM2Q\\archivos sistema_2\\Consolidado\\PRUEBA RAPIDA\\CONTRATISTA LOS MAGNIFICOS S.A.C. - CUNYAS VILA JEAN CARLOS - 03 septiembre,  2020.pdf")
	http.ServeFile(w, r, fp)
}
