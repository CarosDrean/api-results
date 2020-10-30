package models

import jwt "github.com/dgrijalva/jwt-go"

type Claim struct {
	UserResult `json:"user"`
	jwt.StandardClaims
}

