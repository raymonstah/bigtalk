package poster

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raymonstah/bigtalk/domain/poster"
	"github.com/raymonstah/bigtalk/domain/poster/twitter"
	"golang.org/x/xerrors"
	"os"
)

type Event struct {
	Question string `json:"question"`
}
func handleRequest(ctx context.Context, input Event) error {
	p := createPoster(ctx)
	err := p.Post([]byte(input.Question))
	if err != nil {
		return xerrors.Errorf("error posting: %w", err)
	}

	return nil
}

func createPoster(ctx context.Context) poster.Poster {

	twitterKey := os.Getenv("TWITTER_KEY")
	twitterSecret := os.Getenv("TWITTER_SECRET")
	return twitter.New(ctx, twitterKey, twitterSecret)
}

func main() {
	lambda.Start(handleRequest)
}
