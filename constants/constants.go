package constants

type State string

const (
	NotFound           State = "Not found User"
	ErrorUP            State = "Error Updating patient"
	NotFoundMail       State = "Not found mail"
	Accept             State = "Accept"
	InvalidCredentials State = "Invalid Credentials"
	PasswordUpdate     State = "Password Updated"

	IdPruebaRapida       string = "N007-ME000000491"
	IdResultPruebaRapida string = "N007-MF000004612"
)
