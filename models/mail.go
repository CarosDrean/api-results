package models

type Mail struct {
	From     string `json:"from"`
	User     string `json:"user"`
	Password string `json:"password"`
}
