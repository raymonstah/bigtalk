package dao

import (
	"context"
	"time"

	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/lambdas/poller/question"
	"github.com/segmentio/ksuid"
	"golang.org/x/xerrors"
)

// New returns a new DAO that can perform actions on the Questions entity
func New(table dynamo.Table) *DAO {
	return &DAO{
		questionTable: table,
	}
}

// Question is a database representation of a question
type Question struct {
	QuestionID   string `dynamo:"question_id,hash"`
	PostCount    int    `dynamo:"post_count"`
	Question     string `dynamo:"question"`
	CreatedAt    int64  `dynamo:"created_at"`
	LastPostedAt int64  `dynamo:"last_posted_at"`
}

// DAO lets us interact with the Question table
type DAO struct {
	questionTable dynamo.Table
}

// Get gets a question by primary key in dynamo
func (d *DAO) Get(ctx context.Context, questionID string) (question.Question, error) {
	var result Question
	err := d.questionTable.
		Get("question_id", questionID).
		Consistent(true).
		OneWithContext(ctx, &result)
	if err != nil {
		return question.Question{}, xerrors.Errorf("unable to get question by question id ('%v'): %w", questionID, err)
	}
	return transform(result), nil
}

// Poll a question that hasn't been yet polled, or a random one of the next lowest count
func (d *DAO) Poll(ctx context.Context) (question.Question, error) {
	return question.Question{}, nil
}

// Create a new question
func (d *DAO) Create(ctx context.Context, input question.CreateQuestionInput) (question.Question, error) {
	now := time.Now().Unix()

	q := Question{
		QuestionID:   ksuid.New().String(),
		PostCount:    0,
		Question:     input.Question,
		CreatedAt:    now,
		LastPostedAt: now,
	}

	err := d.questionTable.Put(q).RunWithContext(ctx)
	if err != nil {
		return question.Question{}, xerrors.Errorf("unable to create question: %w", err)
	}
	return transform(q), nil
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
