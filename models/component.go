package models

type Component struct {
	ID         string `json:"_id"`
	Name       string `json:"name"`
	CategoryID int    `json:"categoryId"`
	IsDeleted  int
}
