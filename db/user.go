package db

import (
	"database/sql"
	"github.com/CarosDrean/api-results.git/models"
	"log"
	"math/rand"
)

func init(){
	var err error

	for _, sc := range prepStmtsUser{
		sc.stmt, err = DB.Prepare(sc.q)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetUser(id int) []models.User {
	res := []models.User{}
	var item models.User
	// Obtenemos y ejecutamos el get prepared statement.
	get := prepStmtsUser["get"].stmt
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

func GetUsers() []models.User {
	res := []models.User{}
	list := prepStmtsUser["list"].stmt
	rows, err := list.Query()
	if err != nil {
		log.Printf("user: error getting users. err: %v\n", err)
	}
	defer rows.Close()

	// Procesamos los rows.
	for rows.Next() {
		var item models.User
		if err := rows.Scan(&item.ID, &item.Username, &item.Password); err != nil {
			log.Printf("user: error scanning row: %v\n", err)
			continue
		}
		res = append(res, item)
	}
	// Verificamos si hubo error procesando los rows.
	if err := rows.Err(); err != nil {
		log.Printf("user: error reading rows: %v\n", err)
	}

	return res
}

func CreateUser(item models.User) []models.User {
	// Generamos ID único para el nuevo post.
	item.ID = rand.Intn(1000)
	for {
		l := GetUser(item.ID)
		if len(l) == 0 {
			break
		}
		item.ID = rand.Intn(1000)
	}

	// Obtenemos y ejecutamos insert prepared statement.
	insert := prepStmtsUser["insert"].stmt
	_, err := insert.Exec(item.ID, item.Username, item.Password)
	if err != nil {
		log.Printf("user: error inserting user %d into DB: %v\n", item.ID, err)
	}
	return []models.User{item}
}

func UpdateUser(item models.User) {
	// Obtenemos y ejecutamos update prepared statement.
	update := prepStmtsUser["update"].stmt
	_, err := update.Exec(item.ID, item.Username, item.Password)
	if err != nil {
		log.Printf("user: error updating user %d into DB: %v\n", item.ID, err)
	}
}

func DeleteUser(id int) {
	// Obtenemos y ejecutamos delete prepared statement.
	del := prepStmtsUser["delete"].stmt
	_, err := del.Exec(id)
	if err != nil {
		log.Printf("user: error deleting user %d into DB: %v\n", id, err)
	}
}