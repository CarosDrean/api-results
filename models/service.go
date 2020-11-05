package models

type Service struct {
	ID              string
	PersonID        string
	ProtocolID      string
	ServiceDate     string
	ServiceStatusId int
	IsDeleted       int
}
