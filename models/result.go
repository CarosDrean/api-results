package models

type Result struct {
	ID           string `json:"_id"`
	IdService    string `json:"idservide"`
	ProtocolName string `json:"protocolname"`
	Business     string `json:"business"`
	Exam         string `json:"exam"`
	ServiceDate  string `json:"servicedate"`
}
