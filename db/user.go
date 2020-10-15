package db

import (
	"database/sql"
	"github.com/CarosDrean/api-results.git/models"
	"log"
	"math/rand"
)

func GetUser(id int) []models.User {
	res := make([]models.User, 0)
	var item models.User
	get := PrepStmtsUser["get"].Stmt
	err := get.QueryRow(id).Scan(&item.ID, &item.Username, &item.Password)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("user: error getting post. Id: %d, err: %v\n", id, err)
		}
	} else {
		res = append(res, item)
	}
	return res
}

func CreateUser(item models.User) []models.User {
	item.ID = rand.Intn(1000)
	for {
		l := GetUser(item.ID)
		if len(l) == 0 {
			break
		}
		item.ID = rand.Intn(1000)
	}

	insert := PrepStmtsUser["insert"].Stmt
	_, err := insert.Exec(item.ID, item.Username, item.Password)
	if err != nil {
		log.Printf("user: error inserting user %d into DB: %v\n", item.ID, err)
	}
	return []models.User{item}
}

func UpdateUser(item models.User) {
	update := PrepStmtsUser["update"].Stmt
	_, err := update.Exec(item.ID, item.Username, item.Password)
	if err != nil {
		log.Printf("user: error updating user %d into DB: %v\n", item.ID, err)
	}
}
