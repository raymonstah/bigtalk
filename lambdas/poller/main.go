package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/lambdas/poller/question"
	dao "github.com/raymonstah/bigtalk/lambdas/poller/question/ddb"
)

func handleRequest(ctx context.Context) (string, error) {
	session := session.Must(session.NewSession(aws.NewConfig()))

	poller := createQuestionPoller(session, "questions") // todo: use env variable
	q, err := poller.Poll(ctx)
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
	lambda.Start(handleRequest)
}
