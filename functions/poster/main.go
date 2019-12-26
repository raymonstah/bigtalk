package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raymonstah/bigtalk/domain/poster"
	"github.com/raymonstah/bigtalk/domain/poster/twitter"
	"github.com/urfave/cli"
	"log"
	"os"
)

type Event struct {
	Question string `json:"question"`
}

type Handler struct {
	poster poster.Poster
}

func (h *Handler) handle(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		question := record.SNS.Message
		err := h.poster.Post(ctx, []byte(question))
		if err != nil {
			return fmt.Errorf("error posting: %w", err)
		}
	}

	return nil
}

var runNow bool
var twitterConf struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
}

func main() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "now",
			Usage:       "run import now from console",
			Destination: &runNow,
		},
		cli.StringFlag{
			Name:        "twitter-consumer-key",
			EnvVar:      "TWITTER_CONSUMER_KEY",
			Required:    true,
			Destination: &twitterConf.consumerKey,
		},
		cli.StringFlag{
			Name:        "twitter-consumer-secret",
			EnvVar:      "TWITTER_CONSUMER_SECRET",
			Required:    true,
			Destination: &twitterConf.consumerSecret,
		},
		cli.StringFlag{
			Name:        "twitter-access-token",
			EnvVar:      "TWITTER_ACCESS_TOKEN",
			Required:    true,
			Destination: &twitterConf.accessToken,
		},
		cli.StringFlag{
			Name:        "twitter-access-secret",
			EnvVar:      "TWITTER_ACCESS_SECRET",
			Required:    true,
			Destination: &twitterConf.accessSecret,
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

func action(_ *cli.Context) error {
	ctx := context.Background()

	p := twitter.New(twitterConf.consumerKey, twitterConf.consumerSecret, twitterConf.accessToken, twitterConf.accessSecret)
	handler := Handler{poster: p}
	if runNow {
		return handler.handle(ctx, events.SNSEvent{Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: "How are you doing today?",
				},
			},
		}})
	}

	lambda.Start(handler.handle)
	return nil
}
