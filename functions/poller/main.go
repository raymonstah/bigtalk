package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	dao "github.com/raymonstah/bigtalk/domain/question/ddb"
)

// Handler contains everything needed to execute this lambda
type Handler struct {
	poller question.Poller
}

// Handle polls for a question and returns it
func (h *Handler) handle(ctx context.Context) (string, error) {
	q, err := h.poller.Poll(ctx)
	if err != nil {
		return "", err
	}
	return q.Question, nil
}

func createQuestionPoller(session *session.Session, questionTableName string) question.Poller {

	db := dynamo.New(session)
	questionTable := db.Table(questionTableName)
	d := dao.New(questionTable)
	return d
}

func main() {

	s := session.Must(session.NewSession(aws.NewConfig()))
	poller := createQuestionPoller(s, "questions")
	handler := Handler{
		poller: poller,
	}
	lambda.Start(handler.handle)
}
