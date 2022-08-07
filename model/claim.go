package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Authorization
	jwt.StandardClaims
}

type AdminKey struct {
}
