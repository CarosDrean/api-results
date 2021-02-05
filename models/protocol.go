package models

type Protocol struct {
	ID                string `json:"_id"`
	Name              string `json:"name"`
	OrganizationID    string `json:"organizationId"`
	EsoType           int    `json:"esoType"`
	LocationID        string `json:"locationId"`
	GroupOccupationId string `json:"groupOccupationId"`
	IsDeleted         int
}
