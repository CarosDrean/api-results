package models

type Organization struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Mail        string `json:"mail"`
	MailContact string `json:"mailContact"`
	MailMedic   string `json:"mailMedic"`
	UrlMedic    bool   `json:"urlMedic"`
	UrlAdmin    bool   `json:"urlAdmin"`
}

type OrganizationForMailCreateUser struct {
	ID       string `json:"_id"`
	Mail     string `json:"mail"`
	TypeUser string `json:"typeUser"`
}


