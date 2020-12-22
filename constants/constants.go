package constants

type State string
type Role string

const (
	NotFound           State = "Not found User"
	ErrorUP            State = "Error Updating patient"
	NotFoundMail       State = "Not found mail"
	Accept             State = "Accept"
	InvalidCredentials State = "Invalid Credentials"
	PasswordUpdate     State = "Password Updated"

	RouteNewPassword   = "newpassword"
	RouteNewSystemUser = "new-systemuser"
	RouteUserLink      = "userlink"

	ClientURL = "https://resultados.holosalud.org/#/"

	IdPruebaRapida       string = "N007-ME000000491"
	IdResultPruebaRapida string = "N007-MF000004612"

	// systemparameter
	IdConsultings string = "116"

	// roles
	RoleTemp    Role = "Temp"
	RolePatient Role = "Patient"
)
