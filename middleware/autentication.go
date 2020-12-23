package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/CarosDrean/api-results.git/constants"
	"github.com/CarosDrean/api-results.git/db"
	"github.com/CarosDrean/api-results.git/helper"
	"github.com/CarosDrean/api-results.git/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("No se pudo leer")
	}

	publicBytes, err := ioutil.ReadFile("./public.rsa.pub")
	if err != nil {
		log.Fatal("No se pudo leer")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo leer")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se pudo leer")
	}
}

func GenerateJWTExternal(item models.External) string {
	claims := models.ClaimExternal{
		External: item,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 5).Unix(),
			Issuer:    "External",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}
	return result
}

func GenerateJWT(userResult models.UserResult) string {
	claims := models.Claim{
		UserResult: userResult,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}
	return result
}

func ValidateToken(w http.ResponseWriter, r *http.Request) string {
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				_, _ = fmt.Fprintln(w, "Su token ha expirado")
				return "Error"
			case jwt.ValidationErrorSignatureInvalid:
				_, _ = fmt.Fprintln(w, "Su firma de token no coincide")
				return "Error"
			default:
				_, _ = fmt.Fprintln(w, "Su token no es valido")
				return "Error"
			}
		default:
			_, _ = fmt.Fprintln(w, "Su token no es valido")
			return "Error"
		}
	}

	if token.Valid {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, "Bienvenido al sistema")
		return "Accept"
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = fmt.Fprintf(w, "Su token no es valido")
	return "Error"
}

func validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, constants.Role) {
	r.Header.Add("Authorization", r.Header.Get("x-token"))
	var claim models.Claim
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &claim, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				_, _ = fmt.Fprintln(w, "Su token ha expirado")
				return nil, ""
			case jwt.ValidationErrorSignatureInvalid:
				_, _ = fmt.Fprintln(w, "Su firma de token no coincide")
				return nil, ""
			default:
				_, _ = fmt.Fprintln(w, "Su token no es valido")
				return nil, ""
			}
		default:
			_, _ = fmt.Fprintln(w, "Su token no es valido error")
			return nil, ""
		}
	}
	return token, claim.Role
}

func CheckSecurityInternalAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, role := validateToken(w, r)
		if token == nil {
			return
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			if role == "Internal Admin" {
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

func CheckSecurity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Authorization", r.Header.Get("x-token"))
		var claim models.Claim
		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &claim, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					_, _ = fmt.Fprintln(w, "Su token ha expirado")
					return
				case jwt.ValidationErrorSignatureInvalid:
					_, _ = fmt.Fprintln(w, "Su firma de token no coincide")
					return
				default:
					_, _ = fmt.Fprintln(w, "Su token no es valido def")
					return
				}
			default:
				log.Println(err)
				_, _ = fmt.Fprintln(w, "Su token no es valido fin def")
				return
			}
		}

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			//fmt.Fprintf(w, "Bienvenido al sistema")
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "Su token no es valido fin")
			return
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request){
	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error al leer el usuario %s\n", err)
		return
	}
	log.Println(user)
	var stateLogin constants.State
	var id string
	isSystemUser := false
	isPatientParticular := false
	if !user.Particular {
		stateLogin, id = patientBusiness(user)

		if stateLogin == constants.NotFound {
			stateLogin, id = db.ValidateSystemUserLogin(user.User, user.Password)
			isSystemUser = true
		}
	} else {
		isPatientParticular = true
		isSystemUser = false
		stateLogin, id = patientParticular(user)
	}


	switch stateLogin {
	case constants.Accept:
		userResult := models.UserResult{ID: id, Role: getRole(0)}
		if isSystemUser {
			systemUser := db.GetSystemUser(id)
			userResult = models.UserResult{ID: id, Role: getRole(systemUser[0].TypeUser)}
		}
		token := GenerateJWT(userResult)
		result := models.ResponseToken{Token: token}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			fmt.Println(w, "Error al generar el json")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonResult)
	case constants.ErrorUP:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Hubo un error!")
		break
	case constants.NotFoundMail:
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "No se encontro su direccion de correo electronico")
		break
	case constants.NotFound:
		w.WriteHeader(http.StatusForbidden)
		if !isSystemUser || isPatientParticular{
			_, _ = fmt.Fprintf(w, "¡No existe Paciente!")
		} else {
			_, _ = fmt.Fprintf(w, "¡No existe Usuario!")
		}
	case constants.InvalidCredentials:
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "¡Contraseña Incorrecta!")
		break
	case constants.PasswordUpdate:
		w.WriteHeader(http.StatusFound)
		_, _ = fmt.Fprintf(w, "Consulte su correo electronico con las nuevas credenciales.")
		break
	}
}

// en las dos funciones siguientes inicializamos la bd dependiendo del caso, tambien lo dejamos en el main por si acaso
func patientBusiness(user models.UserLogin) (constants.State, string){
	db.DB = helper.Get()
	return db.ValidatePatientLogin(user.User, user.Password)
}

func patientParticular(user models.UserLogin) (constants.State, string){
	db.DB = helper.GetAux()
	return db.ValidatePatientLogin(user.User, user.Password)
}

func getRole(typeUser int)constants.Role{
	switch typeUser {
	case 0:
		return constants.RolePatient
	case 1:
		return "Internal Admin"
	case 2:
		return "External Admin"
	case 3:
		return "External Medic"
	case 4:
		return "External Medic"
	default:
		return ""
	}
}
