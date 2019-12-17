package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	dao "github.com/raymonstah/bigtalk/domain/question/ddb"
)

// Handler contains everything needed to execute this lambda
type Handler struct {
	poller           question.Poller
	sqs              sqsiface.SQSAPI
	queueDestination string
}

// Handle polls for a question and returns it
func (h *Handler) handle(ctx context.Context) error {
	q, err := h.poller.Poll(ctx)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			fmt.Println("no questions to poll for...")
			return nil
		}
		return fmt.Errorf("unable to poll for new question: %w", err)
	}
	_, err = h.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(q.Question),
		QueueUrl:    aws.String(h.queueDestination),
	})
	if err != nil {
		return fmt.Errorf("unable to send message to sqs: %w", err)
	}
	return nil
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
	sqsclient := sqs.New(s)
	handler := Handler{
		poller:           poller,
		sqs:              sqsclient,
		queueDestination: os.Getenv("QUESTIONS_QUEUE"),
	}
	lambda.Start(handler.handle)
}
