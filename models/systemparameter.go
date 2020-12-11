package models

type SystemParameter struct {
	GroupID     int
	ParameterID int `json:"parameterId"`
	Value1      string `json:"value"`
}
