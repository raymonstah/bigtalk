package twitter

import (
	"context"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
)

// Twitter is an implemention of the Poster interface
type Twitter struct {
	client *twitter.Client
}

// New creates a new Poster
func New(consumerKey, consumerSecret, accessToken, accessSecret string) Twitter {

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return Twitter{
		client: client,
	}

}

// Post a tweet!
func (t Twitter) Post(ctx context.Context, content []byte) error {
	_, _, err := t.client.Statuses.Update(string(content), nil)
	if err != nil {
		return fmt.Errorf("error posting tweet: %w", err)
	}
	return nil

}
