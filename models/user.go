package models

type User struct {
	ID       int
	Username string
	Password string
}

type UserResult struct {
	ID   int
	Role string
}