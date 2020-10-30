package models

type User struct {
	ID       int
	Username string
	Password string
}

type UserResult struct {
	ID   string
	Role string
}

type UserLogin struct {
	Email    string
	Password string
}
