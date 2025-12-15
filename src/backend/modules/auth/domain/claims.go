package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func (c *Claims) UserID() string {
	return c.Subject
}
