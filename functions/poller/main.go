package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	dao "github.com/raymonstah/bigtalk/domain/question/ddb"
)

// Handler contains everything needed to execute this lambda
type Handler struct {
	poller   question.Poller
	sqs      sqsiface.SQSAPI
	queueArn string
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

	queueName := getQueueName(h.queueArn)
	queueUrlOutput, err := h.sqs.GetQueueUrlWithContext(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return fmt.Errorf("unable to get queue url from queue name %v: %w", queueName, err)
	}

	_, err = h.sqs.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(q.Question),
		QueueUrl:    queueUrlOutput.QueueUrl,
	})
	if err != nil {
		return fmt.Errorf("unable to send message to sqs: %w", err)
	}
	return nil
}

// extracts queue name from queue arn
func getQueueName(arn string) string {
	lastIndexOfColon := strings.LastIndex(arn, ":")
	return arn[lastIndexOfColon+1:]
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
		poller:   poller,
		sqs:      sqsclient,
		queueArn: os.Getenv("QUESTIONS_QUEUE"),
	}
	lambda.Start(handler.handle)
}
