package models

import (
	"encoding/json"
	"io"
)

type Mail struct {
	Email    string `json:"email"`
	User     string `json:"user"`
	Password string `json:"password"`
	Data     string `json:"data"` // ese elemento es para datos adicionales que se requiera
	Business string `json:"business"`
}

type MailFeedback struct {
	Email   string `json:"mail"`
	User    string `json:"user"`
	Type    string `json:"type"`
	Message string `json:"message"`
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

type MailLink struct {
	Email    string `json:"email"`
	URL      string `json:"url"`
	Business string `json:"business"`
}
