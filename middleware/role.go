package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/models"
	"io"
	"net/http"
)

func RoleInternalAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, role := validateToken(w, r, false)
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
		token, role := validateToken(w, r, true)
		if token == nil {
			return
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			if role == constants.Roles.InternalAdmin{
				next(w, r)
			} else if role == constants.Roles.Temp || role == constants.Roles.ExternalAdmin {
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

func validateCreationForTempOrExternalAdmin(body io.ReadCloser) bool {
	var item models.UserPerson

	_ = json.NewDecoder(body).Decode(&item)
	fmt.Println("hey")
	fmt.Println(item)
	if item.TypeUser == constants.CodeRoles.ExternalAdmin || item.TypeUser == constants.CodeRoles.ExternalMedicNoData {
		return true
	}
	return false
}