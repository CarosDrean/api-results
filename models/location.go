package models

type Location struct {
	ID             string `json:"_id"`
	OrganizationID string `json:"idorganization"`
	Name           string `json:"name"`
	IsDeleted      int
}
