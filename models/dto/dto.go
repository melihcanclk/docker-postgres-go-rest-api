package dto

import "github.com/google/uuid"

type FactsDTO struct {
	ID              int          `json:"id"`
	QuestionContent string       `json:"question_content"`
	Answers         []AnswersDTO `json:"answers"`
}

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserUpdateBodyDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginBodyDTO struct {
	Username string `json:"username,omitempty"` //identity can be username or email, user can login with both
	Email    string `json:"email,omitempty"`    //identity can be username or email, user can login with both
	Password string `json:"password" validate:"required"`
}

type AnswersDTO struct {
	ID         int    `json:"id"`
	AnswerText string `json:"answer_text"`
	IsTrue     bool   `json:"is_true"`
}
