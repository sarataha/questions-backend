// repository/question_repository.go
package repository

import (
	"database/sql"

	"github.com/togglhire/backend-homework/internal/model"
)

type SQLiteRepository struct {
	DB *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{DB: db}
}

func (r *SQLiteRepository) GetAllQuestions() ([]model.Question, error) {
	rows, err := r.DB.Query("SELECT id, body FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.Question
	for rows.Next() {
		var question model.Question
		err := rows.Scan(&question.ID, &question.Body)
		if err != nil {
			return nil, err
		}

		// Fetch options for each question
		options, err := r.getOptionsByQuestionID(question.ID)
		if err != nil {
			return nil, err
		}
		question.Options = options

		questions = append(questions, question)
	}

	return questions, nil
}

func (r *SQLiteRepository) CreateQuestion(question model.Question) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Insert question
	result, err := tx.Exec("INSERT INTO questions (body) VALUES (?)", question.Body)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	questionID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert options
	for _, option := range question.Options {
		_, err := tx.Exec("INSERT INTO options (question_id, body, correct) VALUES (?, ?, ?)",
			questionID, option.Body, option.Correct)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return int(questionID), nil
}

func (r *SQLiteRepository) UpdateQuestion(id int, question model.Question) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Update question
	_, err = tx.Exec("UPDATE questions SET body = ? WHERE id = ?", question.Body, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing options
	_, err = tx.Exec("DELETE FROM options WHERE question_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert updated options
	for _, option := range question.Options {
		_, err := tx.Exec("INSERT INTO options (question_id, body, correct) VALUES (?, ?, ?)",
			id, option.Body, option.Correct)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *SQLiteRepository) DeleteQuestion(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Delete question
	_, err = tx.Exec("DELETE FROM questions WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete options
	_, err = tx.Exec("DELETE FROM options WHERE question_id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *SQLiteRepository) getOptionsByQuestionID(questionID int) ([]model.Option, error) {
	rows, err := r.DB.Query("SELECT body, correct FROM options WHERE question_id = ?", questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []model.Option
	for rows.Next() {
		var option model.Option
		err := rows.Scan(&option.Body, &option.Correct)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, nil
}
