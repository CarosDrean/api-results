package models

type User struct {
	ID       int
	Username string
	Password string
}

type UserResult struct {
	ID   string `json:"_id"`
	Role string `json:"role"`
}

type UserLogin struct {
	User     string
	Password string
}
