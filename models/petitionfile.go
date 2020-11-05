package models

type PetitionFile struct {
	Exam        string `json:"exam"`
	ServiceID   string `json:"serviceId"`
	DNI         string `json:"dni"`
	NameComplet string `json:"nameComplet"`
	ServiceDate string `json:"serviceDate"`
}
