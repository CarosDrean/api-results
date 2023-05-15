package models

type PatientFile struct {
	Exam        string `json:"exam"`
	ServiceID   string `json:"serviceId"`
	DNI         string `json:"dni"`
	NameComplet string `json:"nameComplet"`
	ServiceDate string `json:"serviceDate"`
}

type ExcelPetitionMatrizFile struct {
	Ini             string `json:"ini"`
	Fin             string `json:"fin"`
	OrganizationID  string `json:"organizaionIds"`
}

type ExcelMatrizFile struct {
	Ini             string `json:"ini"`
	Fin             string `json:"fin"`
	VServiceid      string `json:"v_ServiceId"`
	PersonName      string `json:"V_NombrePersona"`
	DocNumber       string `json:"v_DocNumber"`
	Bithdate        string `json:"d_Birthdate"`
	EsoName         string `json:"v_Value1"`
	ProtocolName    string `json:"v_Name"`
	ServiceDate     string `json:"d_ServiceDate"`
	PersonOcupation string `json:"v_CurrentOccupation"`
	Aptitude        string `json:"V_Aptitude"`
	Restriction     string `json:"V_Restriction"`
	IsDelete        int
}

type ExcelInterconsultas struct {
	InterconsultaName    string `json:"v_Name"`
	ServiceId    		 string `json:"v_ServiceId"`
	RepositorioDxId   string `json:"v_DiagnosticRepositoryId"`
	IsDelete        	 int
}

type ExcelRestricciones struct {
	RestrictionName    	 string `json:"v_Name"`
	IsDelete        	 int
}

type ExcelRecomendaciones struct {
	RecomendationName    string `json:"v_Name"`
	IsDelete        	 int
}

type ExcelAluraAptitud struct {
	AptitudName    		 string `json:"v_Value1"`
	IsDelete        	 int
}

type ExcelAptitudEspaciosConfinados struct {
	AptitudName    		 string `json:"v_Value1"`
	IsDelete        	 int
}
