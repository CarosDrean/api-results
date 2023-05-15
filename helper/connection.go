package helper

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/utils"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func Get() (*sql.DB, string) {
	config, err := utils.GetConfiguration()

	if err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("server=%s; user id=%s; password=%s; port=%s; database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Database)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Println("Error connect DB!")
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error connect DB!")
		log.Panic(err)
	}

	fmt.Println("Db is connected!")

	return db, config.Database
}

func GetAux() (*sql.DB, string) {
	config, err := utils.GetConfiguration()

	if err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("server=%s; user id=%s; password=%s; port=%s; database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Databaseaux)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Println("Error connect DB!")
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error connect DB!")
		log.Panic(err)
	}

	fmt.Println("Db aux is connected!")

	return db, config.Databaseaux
}
