package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"net/http"
)

func RoleInternalAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, role := validateToken(w, r)
		if token == nil {
			return
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			if role == constants.Roles.InternalAdmin {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Su rol no le da acceso")
				return
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "Su token no es valido")
			return
		}
	}
}

func RoleInternalAdminOrTempOrExternalAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, role := validateToken(w, r)
		if token == nil {
			return
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			if role == constants.Roles.InternalAdmin{
				next(w, r)
			} else if (role == constants.Roles.Temp || role == constants.Roles.ExternalAdmin) && validateCreationForTempOrExternalAdmin(r) {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Su rol no le da acceso")
				return
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "Su token no es valido")
			return
		}
	}
}

func validateCreationForTempOrExternalAdmin(r *http.Request) bool {
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	if item.TypeUser != constants.CodeRoles.ExternalAdmin && item.TypeUser != constants.CodeRoles.ExternalMedicNoData {
		return false
	}
	return true
}
