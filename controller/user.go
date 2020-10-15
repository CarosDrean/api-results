package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func InitU(){
	var err error
	for _, sc := range db.PrepStmtsUser{
		sc.Stmt, err = db.DB.Prepare(sc.Q)
		if err != nil {
			log.Panic(err)
		}
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item models.User
	_ = json.NewDecoder(r.Body).Decode(&item)
	var ic = db.CreateUser(item)
	_ = json.NewEncoder(w).Encode(ic)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	InitU()

	items := db.GetUser(id)
	_ = json.NewEncoder(w).Encode(items)
}