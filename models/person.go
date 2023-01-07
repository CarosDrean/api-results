package models

type Person struct {
	ID             string `json:"_id"`
	DNI            string `json:"dni"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	FirstLastName  string `json:"firstLastname"`
	SecondLastName string `json:"secondLastname"`
	Mail           string `json:"mail"`
	Sex            int    `json:"sex"`
	Birthday       string `json:"birthday"`
	Phone      	   string `json:"phone"`
	Occupation 	   string `json:"ocupation"`
	Doc            int    `json:"doc"`
	IsDeleted  int
}
