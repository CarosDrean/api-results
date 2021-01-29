package models

type Component struct {
	ID         string  `json:"_id"`
	Name       string  `json:"name"`
	CategoryID int     `json:"categoryId"`
	Price      float32 `json:"price"`
	IsDeleted  int
}
