package main

import (
	"context"
	"github.com/aws/aws-sdk-go/awstesting/mock"
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

func TestHandle(t *testing.T) {
	handler := Handler{poller: &mockPoller{}}
	ctx := context.Background()
	q, err := handler.handle(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "How are you doing today?", q)
}


func TestCreateQuestionPoller(t *testing.T) {
	poller := createQuestionPoller(mock.Session, "question")
	assert.NotNil(t, poller)
}