package models

type Service struct {
	ID              string
	PersonID        string
	ProtocolID      string
	ServiceDate     string
	ServiceStatusId int
	IsDeleted       int
}

type ServicePatient struct {
	ID             string `json:"_id"`
	ServiceDate    string `json:"serviceDate"`
	PersonID       string `json:"personId"`
	DNI            string `json:"dni"`
	Name           string `json:"name"`
	FirstLastName  string `json:"firstLastname"`
	SecondLastName string `json:"secondLastname"`
	Mail           string `json:"mail"`
	Sex            int    `json:"sex"`
}
