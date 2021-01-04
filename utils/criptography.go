package utils

import (
	"crypto/md5"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func EncryptMD5(text string) string {
	data := []byte(text)
	hash := md5.Sum(forSystem(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func forSystem(data []byte) []byte {
	res := make([]byte, 0)
	for _, e := range data {
		res = append(res, e)
		res = append(res, 0)
	}
	return res
}

func ComparePassword(hashedPassword string, password string) bool {
	if hashedPassword != EncryptMD5(password) {
		return false
	}
	/*err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}*/
	return true
}

func Encrypt(password string) string {
	passwordByte := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return ""
	}
	return string(hashedPassword)
}
