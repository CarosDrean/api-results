package models

import "github.com/CarosDrean/api-results.git/constants"

type SystemUser struct {
	ID             int64  `json:"_id"`
	PersonID       string `json:"personId"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	TypeUser       int    `json:"typeUser"`
	OrganizationID string `json:"organizationId"`
	IsDelete       int
}

type UserPerson struct {
	ID               int64  `json:"_id"`
	PersonID         string `json:"personId"`
	UserName         string `json:"username"`
	Password         string `json:"password"`
	TypeUser         int    `json:"typeUser"`
	OrganizationID   string `json:"organizationId"`
	Organization     string `json:"organization"`
	DNI              string `json:"dni"`
	Name             string `json:"name"`
	FirstLastName    string `json:"firstLastname"`
	SecondLastName   string `json:"secondLastname"`
	Mail             string `json:"mail"`
	Sex              int    `json:"sex"`
	Birthday         string `json:"birthday"`
	CodeProfessional string `json:"codeProfessional"` // numero de colegiatura, para el caso inicial siempre va a ser de un medico
}

type UserResult struct {
	ID   string         `json:"_id"`
	Role constants.Role `json:"role"`
}

type UserLogin struct {
	User       string
	Password   string
	Particular bool
}
