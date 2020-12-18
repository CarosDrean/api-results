package models

type Service struct {
	ID               string
	PersonID         string
	ProtocolID       string
	ServiceDate      string
	ServiceStatusId  int
	AptitudeStatusId int
	IsDeleted        int
}

type ServicePatient struct {
	ID               string `json:"_id"`
	ServiceDate      string `json:"serviceDate"`
	PersonID         string `json:"personId"`
	ProtocolID       string `json:"protocolId"`
	AptitudeStatusId int    `json:"aptitude"`
	DNI              string `json:"dni"`
	Name             string `json:"name"`
	FirstLastName    string `json:"firstLastname"`
	SecondLastName   string `json:"secondLastname"`
	Mail             string `json:"mail"`
	Sex              int    `json:"sex"`
	Birthday         string `json:"birthday"`
	// only result covid moment
	Result string `json:"result"`
}

type ServicePatientDiseases struct {
	ID               string `json:"_id"`
	ServiceDate      string `json:"serviceDate"`
	PersonID         string `json:"personId"`
	ProtocolID       string `json:"protocolId"`
	OrganizationID   string `json:"organizationId"`
	AptitudeStatusId int    `json:"aptitude"`
	DNI              string `json:"dni"`
	Name             string `json:"name"`
	FirstLastName    string `json:"firstLastname"`
	SecondLastName   string `json:"secondLastname"`
	Mail             string `json:"mail"`
	Sex              int    `json:"sex"`
	Birthday         string `json:"birthday"`
	Disease          string `json:"disease"`
	Component        string `json:"component"`
	Consulting       string `json:"consulting"`
}
