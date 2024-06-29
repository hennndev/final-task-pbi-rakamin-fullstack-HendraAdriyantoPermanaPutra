package app

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
