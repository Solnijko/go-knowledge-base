package auth

import "time"

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Access   []string  `json:"access"`
	Created  time.Time `json:"created"`
}
