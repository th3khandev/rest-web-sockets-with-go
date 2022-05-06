package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
