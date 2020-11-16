package models

type SystemUser struct {
	ID       string
	PersonID string
	UserName string
	Password string
	TypeUser int
}

type UserResult struct {
	ID   string `json:"_id"`
	Role string `json:"role"`
}

type UserLogin struct {
	User     string
	Password string
}
