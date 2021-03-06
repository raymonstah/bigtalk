package dao

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	"github.com/segmentio/ksuid"
	"gopkg.in/go-playground/validator.v9"
)

const (
	postKey = "postKey"
)

var validate = validator.New()

// New returns a new DAO that can perform actions on the Questions entity
func New(table dynamo.Table) *DAO {
	return &DAO{
		questionTable: table,
	}
}

// Question is a database representation of a question
// Its primary key is the question id
// It also has a GSI of the PostKey / PostCount, where the Partition Key is constant
type Question struct {
	QuestionID   string `dynamo:"question_id,hash"`
	PostKey      string `dynamo:"post_key" index:"poll-index,hash"` // always hardcoded to `postKey`
	PostCount    int    `dynamo:"post_count" index:"poll-index,range"`
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
		return question.Question{}, fmt.Errorf("unable to get question by question id ('%v'): %w", questionID, err)
	}
	return transform(result), nil
}

// List all questions
func (d *DAO) List(ctx context.Context) ([]question.Question, error) {
	var results []Question
	err := d.questionTable.
		Scan().
		Consistent(true).
		AllWithContext(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to list all questions %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt < results[j].CreatedAt
	})
	return transformAll(results), nil
}

// Poll a question that hasn't been yet polled, or a random one of the next lowest count
func (d *DAO) Poll(ctx context.Context) (question.Question, error) {
	var result Question
	err := d.questionTable.Get("post_key", postKey).
		Index("poll-index").
		Limit(1).
		One(&result)
	if err != nil {
		return question.Question{}, fmt.Errorf("unable to poll for question: %w", err)
	}

	return transform(result), nil
}

// Use a question by incrementing its post count and last posted at date
func (d *DAO) Use(ctx context.Context, questionID string) error {
	err := d.questionTable.
		Update("question_id", questionID).
		Add("post_count", 1).
		Set("last_posted_at", time.Now().Unix()).
		RunWithContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to use question: %w", err)
	}

	return nil
}

// Create a new question
func (d *DAO) Create(ctx context.Context, input question.CreateInput) (question.Question, error) {

	if err := validate.Struct(input); err != nil {
		return question.Question{}, err
	}
	now := time.Now().UnixNano()

	q := Question{
		QuestionID: ksuid.New().String(),
		PostKey:    postKey,
		PostCount:  0,
		Question:   input.Question,
		CreatedAt:  now,
	}

	err := d.questionTable.Put(q).RunWithContext(ctx)
	if err != nil {
		return question.Question{}, fmt.Errorf("unable to create question: %w", err)
	}
	return transform(q), nil
}

func (d *DAO) Delete(ctx context.Context, id string) error {
	err := d.questionTable.Delete("question_id", id).
		RunWithContext(ctx)
	if err != nil {
		return fmt.Errorf("error deleting question by id %v: %w", id, err)
	}
	return nil
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

func transformAll(input []Question) []question.Question {
	var results []question.Question
	for _, q := range input {
		results = append(results, transform(q))
	}
	return results
}
