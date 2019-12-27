package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/go-playground/validator/v10"
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

func main() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "now",
			Usage:       "run import now from console",
			Destination: &runNow,
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

var validate = validator.New()

func action(_ *cli.Context) error {
	ctx := context.Background()

	// dynamically get credentials
	tc, err := getTwitterCredentialsFromSecretsManager(ctx)
	if err != nil {
		return fmt.Errorf("unable to get twitter credentials: %w", err)
	}

	if err = validate.Struct(tc); err != nil {
		return err
	}
	p := twitter.New(tc.ConsumerKey, tc.ConsumerSecret, tc.AccessToken, tc.AccessSecret)

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

type twitterCredentials struct {
	ConsumerKey    string `json:"bt-twitter-consumer-key" validate:"required"`
	ConsumerSecret string `json:"bt-twitter-consumer-secret" validate:"required"`
	AccessToken    string `json:"bt-twitter-access-token" validate:"required"`
	AccessSecret   string `json:"bt-twitter-access-secret" validate:"required"`
}

func getTwitterCredentialsFromSecretsManager(ctx context.Context) (twitterCredentials, error) {
	var (
		s          = session.Must(session.NewSession())
		ssm        = secretsmanager.New(s)
		secretName = "bt-secrets"
		version    = "AWSCURRENT" // VersionStage defaults to AWSCURRENT if unspecified
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String(version),
	}

	result, err := ssm.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return twitterCredentials{}, fmt.Errorf("unable to get secret value from secrest manager: %w", err)
	}

	var tc twitterCredentials
	if err = json.Unmarshal([]byte(*result.SecretString), &tc); err != nil {
		return tc, fmt.Errorf("unable to unmarshal secret string from secrets manager: %w", err)
	}

	return tc, nil
}
