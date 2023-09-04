package models

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	AnswerText string `json:"answer" gorm:"text; not null; default: null; column:answer"`
	IsTrue     bool   `json:"is_true" gorm:"not null; default: false; column:is_true"`
	QuestionID uint   `json:"question_id" gorm:"not null; column:question_id"`
}
