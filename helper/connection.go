
package helper

import (
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/utils"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func Get() *sql.DB {
	config, err := utils.GetConfiguration()

	if err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("server=%s; user id=%s; password=%s; port=%s; database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Database)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Db is connected!")

	return db
}