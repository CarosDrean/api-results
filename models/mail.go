package models

import (
	"encoding/json"
	"io"
)

type Mail struct {
	From     string `json:"from"`
	User     string `json:"user"`
	Password string `json:"password"`
	Data     string `json:"data"` // ese elemento es para datos adicionales que se requiera
	Business string `json:"business"`
}

type MailFile struct {
	Email           string `json:"email"`
	FilenameUpload  string `json:"filenameUpload"`
	Description     string `json:"description"`
	NameFileSending string `json:"nameFileSendingNoFormat"`
	FormatFile      string `json:"formatFile"`
}

func (m *MailFile) Decode(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(&m)
}