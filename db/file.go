package db

import (
	_ "database/sql"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/query"
)

type FileDB struct{}

func (db FileDB) GetMatrizOnline(ini string, fin string, org string) ([]models.ExcelMatrizFile, error) {
	res := make([]models.ExcelMatrizFile, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getDataFile"].Q, ini, fin, org)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanExcelMatriz(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetInterconsultas(ser string) ([]models.ExcelInterconsultas, error) {
	res := make([]models.ExcelInterconsultas, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getInterconsultas"].Q, ser)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanInterconsultas(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetRestricciones(ser string) ([]models.ExcelRestricciones, error) {
	res := make([]models.ExcelRestricciones, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getRestriccioens"].Q, ser)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanRestricciones(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetRecomendaciones(repDx string, ser string) ([]models.ExcelRecomendaciones, error) {
	res := make([]models.ExcelRecomendaciones, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getRecomendaciones"].Q, repDx, ser)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanRecomendaciones(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetAlturaAptitud(ser string) ([]models.ExcelAluraAptitud, error) {
	res := make([]models.ExcelAluraAptitud, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getAptitudAltura"].Q, ser)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanAlturaAptitud(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetAptitudEspaciosConfi(ser string) ([]models.ExcelAptitudEspaciosConfinados, error) {
	res := make([]models.ExcelAptitudEspaciosConfinados, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getAptitudEspacios"].Q, ser)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanAptitudEC(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) scanExcelMatriz(row models.RowScanner) (models.ExcelMatrizFile, error) {

	var item models.ExcelMatrizFile

	err := row.Scan(
		&item.VServiceid,
		&item.PersonName,
		&item.DocNumber,
		&item.Bithdate,
		&item.EsoName,
		&item.ProtocolName,
		&item.ServiceDate,
		&item.PersonOcupation,
		&item.Aptitude,
		//&item.Restriction,
	)
	if err != nil {
		return models.ExcelMatrizFile{}, err
	}

	return item, nil
}

func (db FileDB) scanInterconsultas(row models.RowScanner) (models.ExcelInterconsultas, error) {

	var item models.ExcelInterconsultas

	err := row.Scan(
		&item.InterconsultaName,
		&item.ServiceId,
		&item.RepositorioDxId,
	)
	if err != nil {
		return models.ExcelInterconsultas{}, err
	}

	return item, nil
}

func (db FileDB) scanRestricciones(row models.RowScanner) (models.ExcelRestricciones, error) {

	var item models.ExcelRestricciones

	err := row.Scan(
		&item.RestrictionName,
	)
	if err != nil {
		return models.ExcelRestricciones{}, err
	}

	return item, nil
}

func (db FileDB) scanRecomendaciones(row models.RowScanner) (models.ExcelRecomendaciones, error) {

	var item models.ExcelRecomendaciones

	err := row.Scan(
		&item.RecomendationName,
	)
	if err != nil {
		return models.ExcelRecomendaciones{}, err
	}

	return item, nil
}

func (db FileDB) scanAlturaAptitud(row models.RowScanner) (models.ExcelAluraAptitud, error) {

	var item models.ExcelAluraAptitud

	err := row.Scan(
		&item.AptitudName,
	)
	if err != nil {
		return models.ExcelAluraAptitud{}, err
	}

	return item, nil
}

func (db FileDB) scanAptitudEC(row models.RowScanner) (models.ExcelAptitudEspaciosConfinados, error) {

	var item models.ExcelAptitudEspaciosConfinados

	err := row.Scan(
		&item.AptitudName,
	)
	if err != nil {
		return models.ExcelAptitudEspaciosConfinados{}, err
	}

	return item, nil
}
