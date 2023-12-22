package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/togglhire/backend-homework/internal/model"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	fmt.Println("Initializing schema...")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	testDB = db

	// Execute your schema initialization script here
	_, err = testDB.Exec(`
        CREATE TABLE IF NOT EXISTS questions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            body TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS options (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            question_id INTEGER,
            body TEXT NOT NULL,
            correct BOOLEAN NOT NULL,
            FOREIGN KEY (question_id) REFERENCES questions(id)
        );
    `)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	testDB.Close()
}

func TestSQLiteRepository_GetAllQuestions(t *testing.T) {
	repo := NewSQLiteRepository(testDB)

	// Insert test data
	_, err := testDB.Exec(`
		INSERT INTO questions (body) VALUES ('Test Question 1');
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 1', true);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 2', false);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 3', false);
	`)
	assert.NoError(t, err)

	questions, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, questions, 1)

	// Add more assertions based on your data model
}

func TestSQLiteRepository_CreateQuestion(t *testing.T) {
	repo := NewSQLiteRepository(testDB)

	// Create a test question
	question := model.Question{
		Body: "Test Question",
		Options: []model.Option{
			{Body: "Option 1", Correct: true},
			{Body: "Option 2", Correct: false},
		},
	}

	// Insert the question into the database
	id, err := repo.CreateQuestion(question)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	// Retrieve the question from the database and assert its properties
	retrievedQuestion, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, retrievedQuestion, 1)
	assert.Equal(t, question.Body, retrievedQuestion[0].Body)
	assert.Len(t, retrievedQuestion[0].Options, 2)

	// Add more assertions based on your data model
}

func TestSQLiteRepository_UpdateQuestion(t *testing.T) {
	repo := NewSQLiteRepository(testDB)

	// Insert a test question
	_, err := testDB.Exec(`
		INSERT INTO questions (body) VALUES ('Test Question 1');
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 1', true);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 2', false);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 3', false);
	`)
	assert.NoError(t, err)

	// Update the question
	updatedQuestion := model.Question{
		Body: "Updated Test Question",
		Options: []model.Option{
			{Body: "Updated Option 1", Correct: false},
			{Body: "Updated Option 2", Correct: true},
		},
	}

	err = repo.UpdateQuestion(1, updatedQuestion)
	assert.NoError(t, err)

	// Retrieve the updated question and assert its properties
	retrievedQuestion, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, retrievedQuestion, 1)
	assert.Equal(t, updatedQuestion.Body, retrievedQuestion[0].Body)
	assert.Len(t, retrievedQuestion[0].Options, 2)

	// Add more assertions based on your data model
}

func TestSQLiteRepository_DeleteQuestion(t *testing.T) {
	repo := NewSQLiteRepository(testDB)

	// Insert a test question
	_, err := testDB.Exec(`
		INSERT INTO questions (body) VALUES ('Test Question 1');
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 1', true);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 2', false);
		INSERT INTO options (question_id, body, correct) VALUES (1, 'Option 3', false);
	`)
	assert.NoError(t, err)

	// Delete the question
	err = repo.DeleteQuestion(1)
	assert.NoError(t, err)

	// Verify that the question is deleted
	retrievedQuestion, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Empty(t, retrievedQuestion)

	// Add more assertions based on your data model
}
