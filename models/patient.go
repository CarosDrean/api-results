package models

type Patient struct {
	ID             string `json:"_id"`
	DNI            string `json:"dni"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	FirstLastName  string `json:"firstLastname"`
	SecondLastName string `json:"secondLastname"`
	Mail           string `json:"mail"`
	Sex            int    `json:"sex"`
	Birthday       string `json:"birthday"`
}
