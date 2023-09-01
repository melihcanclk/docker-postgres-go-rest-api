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
