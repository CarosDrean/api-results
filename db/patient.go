package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/models"
	"github.com/CarosDrean/api-results.git/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func GetPatientFromDNI(dni string) []models.Patient {
	res := make([]models.Patient, 0)
	var item models.Patient

	tsql := fmt.Sprintf(QueryPatient["getDNI"].Q, dni)
	rows, err := DB.Query(tsql)

	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		return res
	}
	for rows.Next(){
		err := rows.Scan(&item.ID, &item.DNI, &item.Password, &item.Name, &item.FirstLastName, &item.SecondLastName, &item.Mail)
		if err != nil {
			log.Println(err)
			return res
		} else{
			res = append(res, item)
			log.Println(item.Name)
		}
	}
	defer rows.Close()
	return res
}

func ValidatePatientLogin(user string, password string) (helper.State, string){
	items := GetPatientFromDNI(user)
	if len(items) > 0 {
		if ValidateInitPassword(password, items[0]){
			if len(items[0].Mail) != 0{
				newPassword := CreateNewPasswordPatient()
				mail := models.Mail{
					From: items[0].Name,
					Data: newPassword,
				}
				_, err := UpdatePasswordPatient(items[0].ID, newPassword)
				if err != nil {
					return helper.ErrorUP, ""
				}
				Sendmail(mail)
				return helper.PasswordUpdate, ""
			}
			return helper.NotFoundMail, ""
		}
		if items[0].Password == password {
			return helper.Accept, items[0].ID
		}
		return helper.InvalidCredentials, ""
	}
	return helper.NotFoundPatient, ""
}

func UpdatePasswordPatient(id string, password string) (int64, error) {
	ctx := context.Background()
	tsql := fmt.Sprintf(QueryPatient["updatePassword"].Q, id)
	result, err := DB.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", password))
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func ValidateInitPassword(password string, patient models.Patient) bool {
	return patient.DNI == password
}

func CreateNewPasswordPatient() string{
	return utils.StringPassword(8)
}

func Sendmail(mail models.Mail){
	data, err := json.Marshal(mail)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post("http", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(body)
}
