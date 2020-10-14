package controller

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.User
	_ = json.NewDecoder(r.Body).Decode(&item)

	var ic = db.CreateUser(item)

	json.NewEncoder(w).Encode(ic)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items := db.GetUsers()

	json.NewEncoder(w).Encode(items)
}