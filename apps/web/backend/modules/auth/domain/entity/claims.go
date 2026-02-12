package entity

import "time"

type Claims struct {
	ExpiresAt time.Time
	IssuedAt  time.Time
	Issuer    string
	Login     string
	Subject   string
}

func (c *Claims) UserID() string {
	return c.Subject
}
