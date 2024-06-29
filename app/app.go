package app

import (
	"github.com/golang-jwt/jwt/v5"
)

var JWT_KEY = []byte("JWT-TOKEN-KEY")

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
