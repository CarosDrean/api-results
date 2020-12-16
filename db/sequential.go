package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"log"
)

func GetNewID(nodeId int, sequentialId int, prefix string) string {
	return fmt.Sprintf("N%03d-%09d%s", nodeId, sequentialId, prefix)
}

func GetNextSequentialId(nodeId int, tableId int) int {
	item := GetSequential(nodeId, tableId)
	if len(item) > 0 {
		item[0].SequentialID = item[0].SequentialID + 1
		_, err := UpdateSequential(item[0])
		if err != nil {
			log.Println(err)
		}
		return item[0].SequentialID
	} else {
		_, err := CreateSequential(nodeId, tableId)
		if err != nil {
			log.Println(err)
		}
		return 0
	}

}

func GetSequential(nodeId int, tableId int) []models.Sequential {
	res := make([]models.Sequential, 0)
	var item models.Sequential

	tsql := fmt.Sprintf(querySequential["get"].Q, nodeId, tableId)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.NodeID, &item.TableID, &item.SequentialID)
		if err != nil {
			log.Println(err)
		} else {
			res = append(res, item)
		}
	}
	defer rows.Close()
	return res
}

func CreateSequential(nodeId int, tableId int) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(querySequential["insert"].Q)
	fmt.Println(tsql)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("i_NodeId", nodeId),
		sql.Named("i_TableId", tableId),
		sql.Named("i_SecuentialId", 0))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func UpdateSequential(item models.Sequential)(int64, error){
	ctx := context.Background()
	tsql := fmt.Sprintf(querySequential["update"].Q, item.NodeID, item.TableID)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("i_NodeId", item.NodeID),
		sql.Named("i_TableId", item.TableID),
		sql.Named("i_SecuentialId", item.SequentialID))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
