package helper

type State string

const (
	NotFoundPatient    State = "Not found Patient"
	ErrorUP            State = "Error Updating patient"
	NotFoundMail       State = "Not found mail"
	Accept             State = "Accept"
	InvalidCredentials State = "Invalid Credentials"
	PasswordUpdate     State = "Password Updated"

)
