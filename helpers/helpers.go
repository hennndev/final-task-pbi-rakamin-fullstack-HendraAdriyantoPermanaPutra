package helpers

import (
	app "final-task-pbi-fullstackdev/app/jwt"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWT_KEY = []byte("JWT-TOKEN-KEY")

// membuat dan melakukan hashing password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// membuat return value error pada validasi user
func DisplayValidationErrors(field, tag string) string {
	switch tag {
	case "required":
		return field + " tidak boleh kosong"
	case "email":
		return "Email tidak valid"
	case "min":
		return field + " terlalu pendek"
	default:
		return field + " tidak valid"
	}
}

// generate jwt baru saat melakukan login
func GenerateJWT(email string) string {
	expTime := time.Now().Add(time.Hour * 1)
	claims := &app.JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth-jwt",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWT_KEY)

	if err != nil {
		fmt.Println(err)
	}
	return token
}
