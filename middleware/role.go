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
			if role == constants.RoleInternalAdmin {
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

func RoleInternalAdminOrTemp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, role := validateToken(w, r)
		if token == nil {
			return
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			if role == constants.RoleInternalAdmin{
				next(w, r)
			} else if role == constants.RoleTemp && validateCreationForTemp(r) {
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

func validateCreationForTemp(r *http.Request) bool {
	var item models.UserPerson
	_ = json.NewDecoder(r.Body).Decode(&item)
	if item.TypeUser != constants.CodeRoleExternalAdmin && item.TypeUser != constants.CodeRoleExternalMedic {
		return false
	}
	return true
}
