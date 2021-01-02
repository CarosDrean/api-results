package models

import (
	"github.com/CarosDrean/api-results.git/constants"
	jwt "github.com/dgrijalva/jwt-go"
)

type Claim struct {
	ClaimResult `json:"result"`
	jwt.StandardClaims
}

type ClaimResult struct {
	ID   string         `json:"_id"`
	Role constants.Role `json:"role"`
	Data string         `json:"data"`
}
