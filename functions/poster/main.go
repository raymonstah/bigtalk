package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raymonstah/bigtalk/domain/poster"
	"github.com/raymonstah/bigtalk/domain/poster/twitter"
	"golang.org/x/oauth2/clientcredentials"
	"fmt"
	"os"
)

type Event struct {
	Question string `json:"question"`
}

type Handler struct {
	poster poster.Poster
}

func (h *Handler) handle(ctx context.Context, input Event) error {
	err := h.poster.Post(ctx, []byte(input.Question))
	if err != nil {
		return fmt.Errorf("error posting: %w", err)
	}
	return nil
}

func createPoster(ctx context.Context, key, secret string) poster.Poster {

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(ctx)

	return twitter.New(httpClient)

}

func main() {

	ctx := context.Background()
	twitterKey := os.Getenv("TWITTER_KEY")
	twitterSecret := os.Getenv("TWITTER_SECRET")
	p := createPoster(ctx, twitterKey, twitterSecret)
	handler := Handler{poster: p}
	lambda.Start(handler.handle)
}
