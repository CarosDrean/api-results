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
	Result  string `json:"result"`
	Result2 string `json:"result2"`

	OrganizationID   string `json:"organizationId"`
	Organization     string `json:"organization"`
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
	EsoType          int    `json:"esoType"`
}

type ServicePatientOrganization struct {
	ID               string `json:"_id"`
	ServiceDate      string `json:"serviceDate"`
	PersonID         string `json:"personId"`
	ProtocolID       string `json:"protocolId"`
	OrganizationID   string `json:"organizationId"`
	EsoType          int    `json:"esoType"`
	Organization     string `json:"organization"`
	AptitudeStatusId int    `json:"aptitude"`
	DNI              string `json:"dni"`
	Name             string `json:"name"`
	FirstLastName    string `json:"firstLastname"`
	SecondLastName   string `json:"secondLastname"`
	Mail             string `json:"mail"`
	Sex              int    `json:"sex"`
	Birthday         string `json:"birthday"`
	Phone            string `json:"phone"`
	Result2          string `json:"result2"`
}

type ServicePatientExam struct {
	ID               string      `json:"_id"`
	ServiceDate      string      `json:"serviceDate"`
	ProtocolID       string      `json:"protocolId"`
	OrganizationID   string      `json:"organizationId"`
	LocationID       string      `json:"locationId"`
	Protocol         string      `json:"protocol"`
	PriceProtocol    float32     `json:"priceProtocol"`
	TypeDoc          string      `json:"typeDoc"`
	Occupation       string      `json:"occupation"`
	EsoType          int         `json:"esoType"`
	Organization     string      `json:"organization"`
	DNI              string      `json:"dni"`
	Name             string      `json:"name"`
	FirstLastName    string      `json:"firstLastname"`
	SecondLastName   string      `json:"secondLastname"`
	Mail             string      `json:"mail"`
	Sex              int         `json:"sex"`
	AptitudeStatusId int         `json:"aptitude"`
	Birthday         string      `json:"birthday"`
	Components       []Component `json:"components"`
}

type ServiceCovid struct {
	Date           string `json:"date"`
	Name           string `json:"name"`
	FirstLastname  string `json:"lastname"`
	SecondLastName string `json:"secondLastname"`
	DocNumber      string `json:"docNumber"`
	BirthDate      string `json:"birthdate"`
	Age            int    `json:"age"`
	Group          string `json:"group"`
	Occupation     string `json:"occupation"`
	Exam           string `json:"exam"`
	Result         string `json:"result"`
	Sex            int    `json:"sex"`
}
