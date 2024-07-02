package models

type PatientFile struct {
	Exam        string `json:"exam"`
	ServiceID   string `json:"serviceId"`
	DNI         string `json:"dni"`
	NameComplet string `json:"nameComplet"`
	ServiceDate string `json:"serviceDate"`
}

type ExcelPetitionMatrizFile struct {
	Ini            string `json:"ini"`
	Fin            string `json:"fin"`
	OrganizationID string `json:"organizaionIds"`
}

type ExcelMatrizFile struct {
	Ini             string `json:"ini"`
	Fin             string `json:"fin"`
	VPersonId       string `json:"v_PersonId"`
	VServiceid      string `json:"v_ServiceId"`
	PersonName      string `json:"V_NombrePersona"`
	Ape1            string `json:"v_FirstLastName"`
	Ape2            string `json:"v_SecondLastName"`
	Name            string `json:"v_FirstName"`
	DocNumber       string `json:"v_DocNumber"`
	SexType         string `json:"i_SexTypeId"`
	Birthplace      string `json:"v_BirthPlace"`
	Direccion       string `json:"v_AdressLocation"`
	Bithdate        string `json:"d_Birthdate"`
	EsoName         string `json:"v_Value1"`
	OrgName         string `json:"v_OrgName"`
	ExpirationDate  string `json:"d_GlobalExpirationDate"`
	Location        string `json:"v_Location"`
	ProtocolName    string `json:"v_Name"`
	ServiceDate     string `json:"d_ServiceDate"`
	PersonOcupation string `json:"v_CurrentOccupation"`
	Aptitude        string `json:"V_Aptitude"`
	Restriction     string `json:"V_Restriction"`
	IsDelete        int
}

type ExcelInterconsultas struct {
	InterconsultaName string `json:"v_Name"`
	ServiceId         string `json:"v_ServiceId"`
	RepositorioDxId   string `json:"v_DiagnosticRepositoryId"`
	IsDelete          int
}

type ExcelRestricciones struct {
	RestrictionName string `json:"v_Name"`
	IsDelete        int
}

type ExcelRecomendaciones struct {
	RecomendationName string `json:"v_Name"`
	IsDelete          int
}

type ExcelAluraAptitud struct {
	AptitudName string `json:"v_Value1"`
	IsDelete    int
}

type ExcelAptitudEspaciosConfinados struct {
	AptitudName string `json:"v_Value1"`
	IsDelete    int
}

type CustormersValueV1 struct {
	Value1   string `json:"v_Value1"`
	IsDelete int
}

type CustormersValueV2 struct {
	Value1 string `json:"v_Value1"`
}

type ValueFromParameterV1 struct {
	Value1 string `json:"v_Value1"`
}

type DxSingle struct {
	Name string `json:"v_Name"`
}

type CheckDx struct {
	Name string `json:"v_Name"`
}

type NoxiousHabits struct {
	TypeHabitsId string `json:"i_TypeHabitsId"`
	Name         string `json:"v_Value1"`
	Frequency    string `json:"v_Frequency"`
	Comment      string `json:"v_Comment"`
}

type AntecedentesPersonales struct {
	DxDetail string `json:"v_DiagnosticDetail"`
}

type CheckAntePer struct {
	DxDetail string `json:"v_DiagnosticDetail"`
}
