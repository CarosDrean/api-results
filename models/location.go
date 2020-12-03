package models

type Location struct {
	ID             string `json:"_id"`
	OrganizationID string `json:"idOrganization"`
	Name           string `json:"name"`
	IsDeleted      int
}
