package utils

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/models"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func StringPassword(length int) string {
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