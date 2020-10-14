
package helper

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/CarosDrean/api-results/utils"
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

	return db
}