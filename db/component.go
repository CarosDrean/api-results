package db

import (
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetComponentsCategoryId(idCategory string) []models.Component{
	res := make([]models.Component, 0)
	var item models.Component

	tsql := fmt.Sprintf(queryComponent["getCategory"].Q, idCategory)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.IsDeleted)
		if err != nil {
			log.Println(err)
		} else if item.IsDeleted != 1{
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}
