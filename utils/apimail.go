package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"io/ioutil"
	"log"
	"net/http"
)

func SendMail(mail models.Mail, route string) error{
	data, err := json.Marshal(mail)
	if err != nil {
		fmt.Println(err)
	}
	token := loginApiMail()

	req, err := http.NewRequest("POST", constants.ApiMail+ "/" + route, bytes.NewBuffer(data))
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
		log.Println(err)
		return err
	}
	byt := []byte(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(dat)
	return nil
}

func loginApiMail() string{
	secret, err := json.Marshal(map[string]string{
		"secret": constants.SecretApiMail,
	})
	if err != nil {
		fmt.Println(err)
	}
	respToken, err := http.Post(constants.ApiMail+ "/login", "application/json", bytes.NewBuffer(secret))
	if err != nil {
		log.Panic(err)
	}
	defer respToken.Body.Close()
	body, err := ioutil.ReadAll(respToken.Body)
	if err != nil {
		log.Panic(err)
	}
	byt := []byte(string(body))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	return dat["token"].(string)
}
