package models

type SystemUser struct {
	ID             string `json:"_id"`
	PersonID       string `json:"personId"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	TypeUser       int    `json:"typeUser"`
	OrganizationID string `json:"organizationId"`
	IsDelete       int
}

type UserPerson struct {
	ID             string `json:"_id"`
	PersonID       string `json:"personId"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	TypeUser       int    `json:"typeUser"`
	OrganizationID string `json:"organizationId"`
	DNI            string `json:"dni"`
	Name           string `json:"name"`
	FirstLastName  string `json:"firstLastname"`
	SecondLastName string `json:"secondLastname"`
	Mail           string `json:"mail"`
}

type UserResult struct {
	ID   string `json:"_id"`
	Role string `json:"role"`
}

type UserLogin struct {
	User       string
	Password   string
	Particular bool
}
