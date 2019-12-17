package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	"github.com/tj/assert"
	"testing"
)

func withDAO(t *testing.T, callback func(ctx context.Context, dao *DAO)) {
	sesh := session.Must(session.NewSession(aws.NewConfig().
		WithRegion("us-west-2").
		WithEndpoint("http://localhost:8000")))

	// Use sesh
	db := dynamo.New(sesh)
	ctx := context.Background()

	// Create table
	tableName := "question-blah"
	err := db.CreateTable(tableName, Question{}).OnDemand(true).RunWithContext(ctx)
	assert.Nil(t, err)
	// Get table
	table := db.Table(tableName)
	desc, err := table.Describe().RunWithContext(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "question_id", desc.HashKey)
	assert.Empty(t, desc.RangeKey)
	// Delete table when done
	defer func() {
		err := table.DeleteTable().RunWithContext(ctx)
		assert.Empty(t, err)
	}()

	dao := New(table)

	callback(ctx, dao)
}
func TestQuestion(t *testing.T) {

	withDAO(t, func(ctx context.Context, dao *DAO) {

		input := question.CreateQuestionInput{
			Question: "How are you doing today?",
		}
		// Create the question
		q, err := dao.Create(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, input.Question, q.Question)
		assert.NotNil(t, q.QuestionID)

		// Find the question we just created
		got, err := dao.Get(ctx, q.QuestionID)
		assert.Nil(t, err)
		assert.Equal(t, q.QuestionID, got.QuestionID)
	})
}

func TestPoller(t *testing.T) {

	withDAO(t, func(ctx context.Context, dao *DAO) {

		input1 := question.CreateQuestionInput{
			Question: "Question 1",
		}
		// Create the question
		q1, err := dao.Create(ctx, input1)
		assert.Nil(t, err)

		input2 := question.CreateQuestionInput{
			Question: "Question 2",
		}
		// Create the question
		q2, err := dao.Create(ctx, input2)
		assert.Nil(t, err)

		// Use the first question
		err = dao.Use(ctx, q1.QuestionID)
		assert.Nil(t, err)

		// Expect the second question since the first question has already been "used"
		got, err := dao.Poll(ctx)
		assert.Nil(t, err)
		assert.Equal(t, q2.QuestionID, got.QuestionID)
	})
}

func TestPoller_NoQuestions(t *testing.T) {
	withDAO(t, func(ctx context.Context, dao *DAO) {
		_, err := dao.Poll(ctx)
		assert.NotNil(t, err)
		isErrNotFound := errors.Is(err, dynamo.ErrNotFound)
		assert.True(t, isErrNotFound)

	})
}

func TestDAO_List(t *testing.T) {
	withDAO(t, func(ctx context.Context, dao *DAO) {

		input1 := question.CreateQuestionInput{
			Question: "Question 1",
		}
		// Create the question
		q1, err := dao.Create(ctx, input1)
		assert.Nil(t, err)

		input2 := question.CreateQuestionInput{
			Question: "Question 2",
		}
		// Create the question
		q2, err := dao.Create(ctx, input2)
		assert.Nil(t, err)

		allquestions, err := dao.List(ctx)
		fmt.Println(allquestions)
		assert.Nil(t, err)
		assert.Len(t, allquestions, 2)
		assert.Equal(t, q1.QuestionID, allquestions[0].QuestionID)
		assert.Equal(t, q2.QuestionID, allquestions[1].QuestionID)
	})
}
