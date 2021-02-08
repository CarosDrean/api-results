package middleware

import (
	"crypto/rsa"
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
	"strings"
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

func GenerateJWTExternal(item models.ClaimResult) string {
	claims := models.Claim{
		ClaimResult: item,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 3).Unix(),
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

func GenerateJWT(userResult models.ClaimResult) string {
	claims := models.Claim{
		ClaimResult: userResult,
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

func validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, constants.Role, string) {
	r.Header.Add("Authorization", r.Header.Get("x-token"))

	var claim models.Claim
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &claim, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	role := claim.Role

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				_, _ = fmt.Fprintln(w, "Su token ha expirado")
				return nil, "", ""
			case jwt.ValidationErrorSignatureInvalid:
				_, _ = fmt.Fprintln(w, "Su firma de token no coincide")
				return nil, "", ""
			default:
				_, _ = fmt.Fprintln(w, "Su token no es valido")
				return nil, "", ""
			}
		default:
			_, _ = fmt.Fprintln(w, "Su token no es valido error")
			return nil, "", ""
		}
	}
	return token, role, claim.NameDB
}

func reconnectDBParticular(nameDB string) {
	instance := db.DB
	strDB := fmt.Sprint(instance)[1:]
	contains := strings.Contains(strDB, nameDB)
	// corregir esto
	if !contains {
		if nameDB == "HoloCovid" {
			db.DB, _ = helper.GetAux()
		} else {
			db.DB, _ = helper.Get()
		}
	}
}

func CheckSecurity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _, nameDB := validateToken(w, r)
		if token == nil {
			return
		}
		reconnectDBParticular(nameDB)

		if token.Valid {
			w.WriteHeader(http.StatusAccepted)
			//fmt.Fprintf(w, "Bienvenido al sistema")
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "Su token no es valido")
			return
		}
	}
}
