package dao

import (
	"context"
	"time"

	"github.com/raymonstah/bigtalk/lambdas/poller/question"
)

// New returns a new DAO that can perform actions on the Questions entity
func New(questionTableName string) *DAO {
	return &DAO{
		TableName: questionTableName,
	}
}

// Question is a database representation of a question
type Question struct {
	QuestionID   string    `dynamo:"question_id"`
	PostCount    int       `dynamo:"post_count"`
	Question     string    `dynamo:"question"`
	CreatedAt    time.Time `dynamo:"created_at"`
	LastPostedAt time.Time `dynamo:"last_posted_at"`
}

// DAO lets us interact with the Question table
type DAO struct {
	TableName string
}

// Get gets a question by primary key in dynamo
func (d *DAO) Get(ctx context.Context, questionID string) (question.Question, error) {
	return question.Question{}, nil
}

// Poll a question that hasn't been yet polled, or a random one of the next lowest count
func (d *DAO) Poll(ctx context.Context) (question.Question, error) {
	return question.Question{}, nil
}

// Create a new question
func (d *DAO) Create(ctx context.Context, input question.CreateQuestionInput) (question.Question, error) {
	return question.Question{}, nil
}

// a little bit of duplication is better than the wrong abstraction
func transform(input Question) question.Question {
	return question.Question{
		QuestionID:   input.QuestionID,
		Question:     input.Question,
		PostCount:    input.PostCount,
		CreatedAt:    input.CreatedAt,
		LastPostedAt: input.LastPostedAt,
	}
}
