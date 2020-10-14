package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/CarosDrean/api-results/models"
)

type Errorresponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	fmt.Fprintln(w, "Hubo un error")
	// log.Panic(err.Error())
	fmt.Fprintln(w, err.Error())
	var response = Errorresponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
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