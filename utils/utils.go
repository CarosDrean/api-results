package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func CreateNewPassword() string{
	return stringPassword(8)
}

func Sendmail(mail models.Mail){
	data, err := json.Marshal(mail)
	if err != nil {
		fmt.Println(err)
	}
	token := loginApiMail()

	req, err := http.NewRequest("POST", constants.ApiMail+ "/newpassword", bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Post(helper.ApiMail + "/newpassword", "application/json", bytes.NewBuffer(data))
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

func loginApiMail() string{
	secret, err := json.Marshal(map[string]string{
		"secret": constants.SecretApiMail,
	})
	if err != nil {
		fmt.Println(err)
	}
	log.Println(secret)
	respToken, err := http.Post(constants.ApiMail+ "/login", "application/json", bytes.NewBuffer(secret))
	if err != nil {
		log.Panic(err)
	}
	defer respToken.Body.Close()
	body, err := ioutil.ReadAll(respToken.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(body))
	byt := []byte(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat["token"])
	return dat["token"].(string)
}


func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func stringPassword(length int) string {
	return StringWithCharset(length, charset)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	_, _ = fmt.Fprintln(w, "Hubo un error")
	// log.Panic(err.Error())
	_, _ = fmt.Fprintln(w, err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(message)
}

func GetConfiguration() (models.Configuration, error) {
	config := models.Configuration{}
	file, err := os.Open("./configuration.json")

	if err != nil {
		return config, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return config, err
	}

	return config, nil
}