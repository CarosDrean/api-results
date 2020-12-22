package models

type Mail struct {
	From     string `json:"from"`
	User     string `json:"user"`
	Password string `json:"password"`
	Data     string `json:"data"` // ese elemento es para datos adicionales que se requiera
}
