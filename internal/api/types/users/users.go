package users

import "time"

type FullUser struct {
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
type User struct {
	Id                string `json:"id" db:"id"`
	FirstName         string `json:"firstName" db:"first_name"`
	LastName          string `json:"lastName" db:"last_name"`
	Email             string `json:"email" db:"email"`
	Role              int    `json:"role" db:"role"`
	EmailNotification bool   `json:"emailNotification" db:"email_notification"`
}

type AccountCreationRequest struct {
	User UserCreate `json:"accountDetails"`
}
type UserCreate struct {
	FirstName         string `json:"firstName" db:"first_name"`
	LastName          string `json:"lastName" db:"last_name"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password" db:"password"`
	AccessRequest     int    `json:"accessRequest" db:"role"`
	EmailNotification bool   `json:"emailNotification" db:"email_notification"`
}

type UserUpdate struct {
	Id                string `json:"id" db:"id"`
	FirstName         string `json:"firstName" db:"first_name"`
	LastName          string `json:"lastName" db:"last_name"`
	Email             string `json:"email" db:"email"`
	Role              int    `json:"role" db:"role"`
	EmailNotification bool   `json:"emailNotification" db:"email_notification"`
}

type UserLoginForm struct {
	FormData UserLogin `json:"formData"`
}
type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type FrontendUser struct {
	Id                string    `json:"id" db:"id"`
	FirstName         string    `json:"firstName" db:"first_name"`
	LastName          string    `json:"lastName" db:"last_name"`
	Email             string    `json:"email" db:"email"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	Role              int       `json:"role" db:"role"`
	EmailNotification bool      `json:"emailNotification" db:"email_notification"`
}
