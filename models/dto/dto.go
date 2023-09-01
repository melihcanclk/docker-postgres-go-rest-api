package dto

import "github.com/google/uuid"

type FactsDTO struct {
	ID       int
	Question string
	Answer   string
}

type UserDTO struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type UserUpdateBodyDTO struct {
	Username string
	Email    string
	Password string
}

type UserLoginBodyDTO struct {
	Username string `json:"username,omitempty"` //identity can be username or email, user can login with both
	Email    string `json:"email,omitempty"`    //identity can be username or email, user can login with both
	Password string `json:"password" validate:"required"`
}
