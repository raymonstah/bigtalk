package question

import (
	"context"
)

// Question is a high level entity that describes a question should look like
type Question struct {
	QuestionID   string `json:"question_id"`
	PostCount    int    `json:"post_count"`
	Question     string `json:"question"`
	CreatedAt    int64  `json:"created_at"`
	LastPostedAt int64  `json:"last_posted_at"`
}

// CreateQuestionInput is the input used to create a new Question
type CreateQuestionInput struct {
	Question string
}

// Creator allows for questions to be created
type Creator interface {
	Create(ctx context.Context, input CreateQuestionInput) (Question, error)
}

// Poller is an interface that helps retrieves questions
type Poller interface {
	// Poll a question that hasn't been polled yet
	// If all have been polled, get a random one
	Poll(ctx context.Context) (Question, error)
}
