package models

import (
	"github.com/CarosDrean/api-results.git/constants"
	jwt "github.com/dgrijalva/jwt-go"
)

type Claim struct {
	UserResult `json:"user"`
	jwt.StandardClaims
}

type ClaimExternal struct {
	External `json:"external"`
	jwt.StandardClaims
}

type External struct {
	OrganizationID string         `json:"organizationId"`
	Role           constants.Role `json:"role"`
}
