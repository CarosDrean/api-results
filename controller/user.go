package controller

import (
	"encoding/json"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

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

	items := db.GetUser(id)
	_ = json.NewEncoder(w).Encode(items)
}