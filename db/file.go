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

func (db FileDB) GetValueCustomerV1(ser string, comp string, CompField string) ([]models.CustormersValueV1, error) {
	res := make([]models.CustormersValueV1, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getCustomeValueV1"].Q, ser, comp, CompField)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanValueCustomizeV1(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetValueCustomerV2(ser string, comp string, CompField string) ([]models.CustormersValueV2, error) {
	res := make([]models.CustormersValueV2, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getCustomeValueV2"].Q, ser, comp, CompField)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanValueCustomizeV2(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetValueFromParameterV1(group string, parameter string) ([]models.ValueFromParameterV1, error) {
	res := make([]models.ValueFromParameterV1, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getValueFromParametersV1"].Q, group, parameter)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanValueFromParameterV1(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetDxSingle(service string, component string) ([]models.DxSingle, error) {
	res := make([]models.DxSingle, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getDxSingle"].Q, service, component)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanDxSingle(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetCheckDx(service string, diseases string) ([]models.CheckDx, error) {
	res := make([]models.CheckDx, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getCheckDx"].Q, service, diseases)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanCheckDx(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetNoxiusHabitats(service string, TypeHabitsId string) ([]models.NoxiousHabits, error) {
	res := make([]models.NoxiousHabits, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getNoxiusHabitats"].Q, service, TypeHabitsId)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanNoxiousHabits(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetAntecedentesPersonales(IdPerson string) ([]models.AntecedentesPersonales, error) {
	res := make([]models.AntecedentesPersonales, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getAntecedentesPersonales"].Q, IdPerson)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanAntecedentesPersonales(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) GetCheckAntePer(IdPerson string, Disease string) ([]models.CheckAntePer, error) {
	res := make([]models.CheckAntePer, 0)

	tsql := fmt.Sprintf(query.ExcelFile["getCheckAntePerso"].Q, IdPerson, Disease)

	rows, err := DB.Query(tsql)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		item, _ := db.scanCheckAntePer(rows)

		res = append(res, item)
	}

	return res, nil
}

func (db FileDB) scanExcelMatriz(row models.RowScanner) (models.ExcelMatrizFile, error) {

	var item models.ExcelMatrizFile

	err := row.Scan(
		&item.VPersonId,
		&item.VServiceid,
		&item.PersonName,
		&item.Ape1,
		&item.Ape2,
		&item.Name,
		&item.DocNumber,
		&item.SexType,
		&item.Birthplace,
		&item.Direccion,
		&item.Bithdate,
		&item.EsoName,
		&item.OrgName,
		&item.ExpirationDate,
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

func (db FileDB) scanValueCustomizeV1(row models.RowScanner) (models.CustormersValueV1, error) {

	var item models.CustormersValueV1

	err := row.Scan(
		&item.Value1,
	)
	if err != nil {
		return models.CustormersValueV1{}, err
	}

	return item, nil
}

func (db FileDB) scanValueCustomizeV2(row models.RowScanner) (models.CustormersValueV2, error) {

	var item models.CustormersValueV2

	err := row.Scan(
		&item.Value1,
	)
	if err != nil {
		return models.CustormersValueV2{}, err
	}

	return item, nil
}

func (db FileDB) scanValueFromParameterV1(row models.RowScanner) (models.ValueFromParameterV1, error) {

	var item models.ValueFromParameterV1

	err := row.Scan(
		&item.Value1,
	)
	if err != nil {
		return models.ValueFromParameterV1{}, err
	}

	return item, nil
}

func (db FileDB) scanDxSingle(row models.RowScanner) (models.DxSingle, error) {

	var item models.DxSingle

	err := row.Scan(
		&item.Name,
	)
	if err != nil {
		return models.DxSingle{}, err
	}

	return item, nil
}

func (db FileDB) scanCheckDx(row models.RowScanner) (models.CheckDx, error) {

	var item models.CheckDx

	err := row.Scan(
		&item.Name,
	)
	if err != nil {
		return models.CheckDx{}, err
	}

	return item, nil
}

func (db FileDB) scanNoxiousHabits(row models.RowScanner) (models.NoxiousHabits, error) {

	var item models.NoxiousHabits

	err := row.Scan(
		&item.TypeHabitsId,
		&item.Name,
		&item.Frequency,
		&item.Comment,
	)
	if err != nil {
		return models.NoxiousHabits{}, err
	}

	return item, nil
}

func (db FileDB) scanAntecedentesPersonales(row models.RowScanner) (models.AntecedentesPersonales, error) {

	var item models.AntecedentesPersonales

	err := row.Scan(
		&item.DxDetail,
	)
	if err != nil {
		return models.AntecedentesPersonales{}, err
	}

	return item, nil
}

func (db FileDB) scanCheckAntePer(row models.RowScanner) (models.CheckAntePer, error) {

	var item models.CheckAntePer

	err := row.Scan(
		&item.DxDetail,
	)
	if err != nil {
		return models.CheckAntePer{}, err
	}

	return item, nil
}
