package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/raymonstah/bigtalk/domain/poster"
	"github.com/tj/assert"
	"testing"
	"time"
)

type mockPoster struct {
	poster.Poster
}

func (m *mockPoster) Post(ctx context.Context, content []byte) error {
	return nil
}

func TestHandle(t *testing.T) {
	p := &mockPoster{}
	handler := Handler{poster: p}
	ctx := context.Background()
	err := handler.handle(ctx, events.SNSEvent{Records: []events.SNSEventRecord{
		{
			SNS: events.SNSEntity{
				Signature:         "",
				MessageID:         "",
				Type:              "",
				TopicArn:          "",
				MessageAttributes: nil,
				SignatureVersion:  "",
				Timestamp:         time.Time{},
				SigningCertURL:    "",
				Message:           "",
				UnsubscribeURL:    "",
				Subject:           "",
			},
		},
	}})
	assert.Nil(t, err)
}
