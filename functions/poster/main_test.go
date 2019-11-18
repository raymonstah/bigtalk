package main

import (
	"context"
	"github.com/raymonstah/bigtalk/domain/poster"
	"github.com/tj/assert"
	"testing"
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
	err := handler.handle(ctx, Event{Question: "How are you doing today?"})
	assert.Nil(t, err)
}

func TestCreatePoster(t *testing.T) {
	p := createPoster(context.Background(), "key", "secret")
	assert.NotNil(t, p)
}