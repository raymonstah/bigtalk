package question

import (
	"context"
	"time"
)

// Question is a high level entity that describes a question should look like
type Question struct {
	QuestionID   string
	PostCount    int       `dynamo:"post_count"`
	Question     string    `dynamo:"question"`
	CreatedAt    time.Time `dynamo:"created_at"`
	LastPostedAt time.Time `dynamo:"last_posted_at"`
}

// CreateQuestionInput is the input used to create a new Question
type CreateQuestionInput struct {
	Question string
}

// Poller is an interface that helps retrieves questions
type Poller interface {
	// Get by ID
	Get(ctx context.Context, questionID string) (Question, error)
	// Poll a question that hasn't been polled yet
	// If all have been polled, get a random one
	Poll(ctx context.Context) (Question, error)
}