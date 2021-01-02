package models

type Professional struct {
	PersonID     string
	ProfessionID int // 32 es medico
	Code         string
	IsDeleted    int
}
