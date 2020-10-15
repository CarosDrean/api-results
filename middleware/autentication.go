package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/controller"
	"github.com/CarosDrean/api-results.git/models"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error al leer el usuario %s\n", err)
		return
	}
	users := controller.GetPatientFromLogin(user)
	if len(users) > 0 {
		if users[0].Password == user.Password {
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprintf(w, "Bienvenido al sistema")
			return
			// aqui devolver login
		} else{
			w.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintln(w, "¡Contraseña no valida!")
			// contraseña invalida
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, "¡Usuario no existe!")
		return
		// no se encontraron registros
	}
}