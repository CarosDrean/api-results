package models

type Mail struct {
	From     string `json:"from"`
	User     string `json:"user"`
	Password string `json:"password"`
	Data     string `json:"data"` // ese elemento es para datos adicionales que se requiera
	Business string `json:"business"`
}

type MailFile struct {
	From     string `json:"from"`
	File     string `json:"file"`
	Business string `json:"business"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}
