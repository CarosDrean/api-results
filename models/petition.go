package models

type PetitionProgrammation struct {
	PersonId               string `json:"v_PersonId"`
	DocType                int    `json:"i_DocTypeId"`
	DocNumber              string `json:"v_DocNumber"`
	FirstName              string `json:"v_FirstName"`
	FirstLastName          string `json:"v_FirstLastName"`
	SecondLastName         string `json:"v_SecondLastName"`
	SexTypeId              int    `json:"i_SexTypeId"`
	Birthdate              string `json:"d_Birthdate"`
	TelephoneNumber        string `json:"v_TelephoneNumber"`
	CurrentOccupation      string `json:"v_CurrentOccupation"`
	DateProgramming        string `json:"d_DateProgramming"`
	ServiceTypeId          int    `json:"i_ServiceTypeId"`
	PersonProgramming      string `json:"v_PersonProgramming"`
	ResponsableProgramming string `json:"v_ResponsableProgramming"`
	CalendarId_2           string `json:"v_CalendarId_2"`
	WorkersCondition       string `json:"v_WorkersCondition"`
	FactCR                 string `json:"v_FactCR"`
	NombreProyecto         string `json:"v_NombreProyecto"`
	OrganizationId         string `json:"v_OrganizationId"`
	ProtocolId             string `json:"v_ProtocolId"`
	Deleted                int    `json:"d_deleted"`
	PetitionStatus         int    `json:"v_PetitionStatus"`
	Comentary              string `json:"v_Comentary"`
}

type MailConsultaCardiologica struct {
	Email      string `json:"mail"`
	Dni        string `json:"dni"`
	Nombre     string `json:"nombre"`
	Apepaterno string `json:"apepaterno"`
	Apematerno string `json:"apematerno"`
	Telefono   string `json:"telefono"`
	Direccion  string `json:"direccion"`
	Sexo       string `json:"sexo"`
	Dob        string `json:"dob"`
	Fecha      string `json:"fecha"`
	Mensaje    string `json:"mensaje"`
}
