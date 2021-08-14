package models

import (
	"encoding/json"
	"fmt"
	"io"
)

type Error struct {
	Code      int         `json:"code"`
	ErrorData interface{} `json:"error"`
	Message   string      `json:"message"`
	Where     string      `json:"where"`
}

func (e *Error) Decode(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(&e)
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Error: %v", e.Code, e.Message, e.ErrorData)
}
