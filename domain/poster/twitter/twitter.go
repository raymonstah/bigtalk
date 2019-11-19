package twitter

import (
	"context"
	"github.com/dghubble/go-twitter/twitter"
	"fmt"
	"net/http"
)


// Twitter is an implemention of the Poster interface
type Twitter struct {
	client *twitter.Client
}

// New creates a new Poster
func New(httpclient *http.Client) Twitter {

	// Twitter client
	client := twitter.NewClient(httpclient)

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
