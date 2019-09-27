package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/lambdas/poller/question"
	"github.com/tj/assert"
)

func withDAO(t *testing.T, callback func(ctx context.Context, dao *DAO)) {
	session := session.Must(session.NewSession(aws.NewConfig().
		WithRegion("us-west-2").
		WithEndpoint("http://localhost:8000")))

	// Use session
	db := dynamo.New(session)
	ctx := context.Background()

	// Create table
	tableName := "question-blah"
	err := db.CreateTable(tableName, Question{}).OnDemand(true).RunWithContext(ctx)
	fmt.Println(err)
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
		q, err := dao.Create(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, input.Question, q.Question)
		assert.NotNil(t, q.QuestionID)

		got, err := dao.Get(ctx, q.QuestionID)
		assert.Nil(t, err)
		assert.Equal(t, q.QuestionID, got.QuestionID)

	})

}
