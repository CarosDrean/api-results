package constants

type State string
type Role string
type IdConsultElectro string
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
	RouteUploadFile    = "upload-file"
	RouteSendFile      = "send-file"

	ClientURL = "https://resultados.holosalud.org/#/"

	IdPruebaRapida       string = "N007-ME000000491"
	IdResultPruebaRapida string = "N007-MF000004612"

	// codigo de bd principal  del hisopado
	IdPruebaHisopado       string = "N009-ME000000567"
	// codigo de bd auxiliar del hisopado
	IdPruebaHisopadoAux       string = "N009-ME000000529"
	IdResultPruebaHisopado string = "N009-MF000004572"
	IdCardio               string = "N009-PR000002876"

	// systemparameter
	IdConsultings string = "116"

	CodeAccessClient = 3003
)

var (
	Roles = RolesModel{
		Temp:          "Temp",
		Patient:       "Patient",
		InternalAdmin: "Internal Admin",
		ExternalAdmin: "External Admin",
		ExternalMedic: "External Medic",
		Accounting:    "Accounting",
	}
	CodeRoles = CodeRolesModel{
		Patient:             0,
		InternalAdmin:       1,
		ExternalAdmin:       2,
		ExternalMedic:       3,
		ExternalMedicNoData: 4,
		Accounting:          5,
	}
)

type RolesModel struct {
	Temp          Role
	Patient       Role
	InternalAdmin Role
	ExternalAdmin Role
	ExternalMedic Role
	Accounting    Role
}

type CodeRolesModel struct {
	Patient             int
	InternalAdmin       int
	ExternalAdmin       int
	ExternalMedic       int
	ExternalMedicNoData int
	Accounting          int
}
