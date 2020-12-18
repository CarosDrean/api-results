package models

type Protocol struct {
	ID             string `json:"_id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
	LocationID     string
	IsDeleted      int
}
