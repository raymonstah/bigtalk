package twitter

import (
	"context"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/xerrors"
)

type Poster struct {
	client *twitter.Client
}

func New(ctx context.Context, key, secret string) Poster {
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(ctx)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return Poster{
		client: client,
	}

}

func (t Poster) Post(content []byte) error {
	_, _, err := t.client.Statuses.Update(string(content), nil)
	if err != nil {
		return xerrors.Errorf("error posting tweet: %w", err)
	}
	return nil

}
