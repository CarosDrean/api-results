package models

type Calendar struct {
	ID               string
	PersonID         string
	ServiceID        string
	ProtocolID       string
	CalendarStatusID int
	CircuitStart     string
	CircuitEnd       string
	Date             string
}
