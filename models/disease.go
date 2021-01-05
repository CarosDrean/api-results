package models

type Disease struct {
	ID        string `json:"_id"`
	CIE10ID   string `json:"cie10Id"`
	Name      string `json:"name"`
	IsDeleted int
}
