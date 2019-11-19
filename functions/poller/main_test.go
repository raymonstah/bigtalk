package main

import (
	"context"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/raymonstah/bigtalk/domain/question"
	"github.com/tj/assert"
	"testing"
)

type mockPoller struct {
	question.Poller
}

func (m mockPoller) Poll(ctx context.Context) (question.Question, error) {
	return question.Question{
		QuestionID:   "123",
		PostCount:    0,
		Question:     "How are you doing today?",
		CreatedAt:    0,
		LastPostedAt: 0,
	}, nil
}


type mockSQS struct {sqsiface.SQSAPI}
func (m *mockSQS) SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return nil, nil
}

func TestHandle(t *testing.T) {
	handler := Handler{
		poller:           &mockPoller{},
		sqs:              &mockSQS{},
		queueDestination: "question-queue",
	}
	ctx := context.Background()
	err := handler.handle(ctx)
	assert.Nil(t, err)
}


func TestCreateQuestionPoller(t *testing.T) {
	poller := createQuestionPoller(mock.Session, "question")
	assert.NotNil(t, poller)
}