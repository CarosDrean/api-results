package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"log"
	"net/http"
)

func Init(){
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

	json.NewEncoder(w).Encode(ic)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Init()
	items := db.GetUsers()

	json.NewEncoder(w).Encode(items)
}