package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func InitP(){
	var err error
	for _, sc := range db.PrepStmtsUser{
		sc.Stmt, err = db.DB.Prepare(sc.Q)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	dni, _ := strconv.Atoi(params["dni"])

	InitP()
	items := db.GetPatient(dni)

	_ = json.NewEncoder(w).Encode(items)
}
