package models

import (
	"encoding/json"
	"io"
)

type MailFileResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
	Format  string `json:"format"`
}

func (m *MailFileResponse) Decode(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(&m)
}

type MailResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (mr *MailResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &mr)
}
