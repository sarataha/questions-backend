package repository

import "github.com/togglhire/backend-homework/internal/model"

type QuestionRepository interface {
	GetAllQuestions() ([]model.Question, error)
	CreateQuestion(question model.Question) (int, error)
	UpdateQuestion(id int, question model.Question) error
	DeleteQuestion(id int) error
}
