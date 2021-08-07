package models

type MailFileRes struct {
	Message string `json:"message"`
	Data    string `json:"data"`
	Format  string `json:"format"`
}
