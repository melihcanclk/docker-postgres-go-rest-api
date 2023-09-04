package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	QuestionContent string   `json:"question_content" gorm:"text; not null; default: null; column:question_content"`
	Answers         []Answer `json:"answers" gorm:"cascade; foreignkey:QuestionID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
