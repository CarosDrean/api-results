package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type SequentialDB struct {}

func (db SequentialDB) NewID(nodeId int, sequentialId int, prefix string) string {
	return fmt.Sprintf("N%03d-%s%09d", nodeId, prefix, sequentialId)
}

func (db SequentialDB) NextSequentialId(nodeId int, tableId int) int {
	item, err := db.Get(nodeId, tableId)
	if err != nil {
		checkError(err, "Get", "Sequential DB", "Next")
	}
	if item.TableID == 0 && item.NodeID == 0 { // valida si el objeto item esta vacio
		_, err := db.Create(nodeId, tableId)
		if err != nil {
			checkError(err, "Create", "Sequential DB", "Next")
		}
		return 0
	}
	item.SequentialID = item.SequentialID + 1
	_, err = db.Update(item)
	if err != nil {
		checkError(err, "Update", "Sequential DB", "Next")
	}
	return item.SequentialID

}

func (db SequentialDB) Get(nodeId int, tableId int) (models.Sequential, error) {
	res := make([]models.Sequential, 0)

	tsql := fmt.Sprintf(query.Sequential["get"].Q, nodeId, tableId)
	rows, err := DB.Query(tsql)

	err = db.scan(rows, err, &res, "Sequential DB", "Get")
	if err != nil {
		return models.Sequential{}, err
	}
	if len(res) == 0 {
		return models.Sequential{}, nil
	}
	defer rows.Close()
	return res[0], nil
}

func (db SequentialDB) Create(nodeId int, tableId int) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Sequential["insert"].Q)
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

func (db SequentialDB) Update(item models.Sequential)(int64, error){
	ctx := context.Background()
	tsql := fmt.Sprintf(query.Sequential["update"].Q, item.NodeID, item.TableID)
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

func (db SequentialDB) scan(rows *sql.Rows, err error, res *[]models.Sequential, ctx string, situation string) error {
	var item models.Sequential
	if err != nil {
		checkError(err, situation, ctx, "Reading rows")
		return err
	}
	for rows.Next() {
		err := rows.Scan(&item.NodeID, &item.TableID, &item.SequentialID)
		if err != nil {
			checkError(err, situation, ctx, "Scan rows")
			return err
		} else {
			*res = append(*res, item)
		}
	}
	return nil
}

