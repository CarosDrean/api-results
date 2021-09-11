package models

import "time"

type SystemUser struct {
	ID               int64     `json:"_id"`
	PersonID         string    `json:"personId"`
	UserName         string    `json:"username"`
	Password         string    `json:"password"`
	TypeUser         int       `json:"typeUser"`
	OrganizationID   string    `json:"organizationId"`
	CodeProfessional string    `json:"codeProfessional"` // numero de colegiatura, para el caso inicial siempre va a ser de un medico
	AccessClient     bool      `json:"accessClient"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	IsDelete         int
}

type UserPerson struct {
	ID               int64     `json:"_id"`
	PersonID         string    `json:"personId"`
	UserName         string    `json:"username"`
	Password         string    `json:"password"`
	TypeUser         int       `json:"typeUser"`
	OrganizationID   string    `json:"organizationId"`
	Organization     string    `json:"organization"`
	DNI              string    `json:"dni"`
	Name             string    `json:"name"`
	FirstLastName    string    `json:"firstLastname"`
	SecondLastName   string    `json:"secondLastname"`
	Mail             string    `json:"mail"`
	Sex              int       `json:"sex"`
	Birthday         string    `json:"birthday"`
	AccessClient     bool      `json:"accessClient"`
	CodeProfessional string    `json:"codeProfessional"` // numero de colegiatura, para el caso inicial siempre va a ser de un medico
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type UserLogin struct {
	User       string
	Password   string
	Particular bool
}
