package models

type Result struct {
	ID           string `json:"_id"`
	IdService    string `json:"idservice"`
	ProtocolName string `json:"protocolname"`
	Business     string `json:"business"`
	Exam         string `json:"exam"`
	ServiceDate  string `json:"servicedate"`
	Result       string `json:"result"`
	Result2      string `json:"result2"`
}
