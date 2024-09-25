package users

import "time"

type User struct {
	Id                string    `json:"id" db:"id"`
	FirstName         string    `json:"firstName" db:"first_name"`
	LastName          string    `json:"lastName" db:"last_name"`
	Email             string    `json:"email" db:"email"`
	Password          string    `json:"password" db:"password"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
	Role              int       `json:"role" db:"role"`
	EmailNotification bool      `json:"emailNotification" db:"email_notification"`
}

type UserCreate struct {
	FirstName         string `json:"firstName" db:"first_name"`
	LastName          string `json:"lastName" db:"last_name"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password" db:"password"`
	AccessRequest     int    `json:"accessRequest" db:"role"`
	EmailNotification bool   `json:"emailNotification" db:"email_notification"`
}
