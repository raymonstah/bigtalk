package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
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
	snsAPI   snsiface.SNSAPI
	topicArn string
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

	_, err = h.snsAPI.PublishWithContext(ctx, &sns.PublishInput{
		Message:  aws.String(q.Question),
		TopicArn: aws.String(h.topicArn),
	})
	if err != nil {
		return fmt.Errorf("unable to publish message to sns topic %v: %w", h.topicArn, err)
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
	snsAPI := sns.New(s)
	handler := Handler{
		poller:   poller,
		snsAPI:   snsAPI,
		topicArn: os.Getenv("QUESTIONS_TOPIC"),
	}
	lambda.Start(handler.handle)
}
