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
	RoleTemp          Role = "Temp"
	RolePatient       Role = "Patient"
	RoleInternalAdmin Role = "Internal Admin"
	RoleExternalAdmin Role = "External Admin"
	RoleExternalMedic Role = "External Medic"

	// code role
	CodeRolePatient             = 0
	CodeRoleInternalAdmin       = 1
	CodeRoleExternalAdmin       = 2
	CodeRoleExternalMedic       = 3
	CodeRoleExternalMedicNoData = 4
)
