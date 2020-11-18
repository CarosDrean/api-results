package models

type SystemUser struct {
	ID       string `json:"_id"`
	PersonID string `json:"personId"`
	UserName string `json:"username"`
	Password string `json:"password"`
	TypeUser int    `json:"typeUser"`
}

type UserResult struct {
	ID   string `json:"_id"`
	Role string `json:"role"`
}

type UserLogin struct {
	User     string
	Password string
}
