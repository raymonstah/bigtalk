package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Config struct {
	skipLambda      bool // if the lambda uploading should be skipped
	btLambdasBucket string
	dir             string
	cf              struct {
		dir       string
		stackname string
	}
}

var config Config

func main() {
	app := cli.NewApp()
	app.Usage = "deploy app to cloud"
	app.Version = "latest"
	app.EnableBashCompletion = true
	app.Action = cli.ActionFunc(action)
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "skip-lambda",
			Destination: &config.skipLambda,
		},
		cli.StringFlag{
			Name:        "bt-lambdas-bucket",
			Usage:       "where the s3 bucket is",
			Value:       "bt-lambdas",
			Destination: &config.btLambdasBucket,
		},
		cli.StringFlag{
			Name:        "path",
			Usage:       "path to artifact resources (zip files for lambdas)",
			Required:    true,
			Destination: &config.dir,
		},
		cli.StringFlag{
			Name:        "cf-directory",
			Usage:       "path to cloudformation template",
			Required:    true,
			Destination: &config.cf.dir,
		},
		cli.StringFlag{
			Name:        "cf-stack-name",
			Usage:       "name of the cloudformation stack",
			Value:       "bt-stack",
			Destination: &config.cf.stackname,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR:", err)
		fmt.Println()
		os.Exit(1)
	}
}

func action(_ *cli.Context) error {
	ctx := context.Background()
	var (
		s     = session.Must(session.NewSession())
		s3API = s3.New(s)
		cfAPI = cloudformation.New(s)
	)

	if config.skipLambda {
		log.Fatal("whoops, skip lambda not implemented")
		// todo get lambda ids
	}

	lambdaVersions, err := uploadLambdas(ctx, s3API, config.dir)
	if err != nil {
		return fmt.Errorf("unable to upload lambdas: %w", err)
	}

	err = uploadCloudFormation(ctx, cfAPI, lambdaVersions)
	if err != nil {
		return fmt.Errorf("unable to upload cloudformation stack: %w", err)
	}

	return nil
}

func uploadCloudFormation(ctx context.Context, api cloudformationiface.CloudFormationAPI, versions map[string]string) error {
	bytes, err := ioutil.ReadFile(config.cf.dir)
	if err != nil {
		return fmt.Errorf("error reading file %v: %w", config.cf.dir, err)
	}

	parameters := makeParams(string(bytes), versions)
	_, err = api.UpdateStackWithContext(ctx, &cloudformation.UpdateStackInput{
		Capabilities: []*string{
			aws.String(cloudformation.CapabilityCapabilityNamedIam),
			aws.String(cloudformation.CapabilityCapabilityAutoExpand),
		},
		Parameters:   parameters,
		StackName:    aws.String(config.cf.stackname),
		TemplateBody: aws.String(string(bytes)),
	})
	if err != nil {
		return fmt.Errorf("unable to update stack: %w", err)
	}
	return nil
}

func makeParams(templateBody string, lambdaVersions map[string]string) []*cloudformation.Parameter {
	var parameters []*cloudformation.Parameter

	var keys, _ = extractParameterNames(templateBody)
	for _, key := range keys {
		switch key {
		case "BigTalkBucket":
			parameters = appendParameter(parameters, key, config.btLambdasBucket)
		default:
			for k, v := range lambdaVersions {
				if k == strings.ToLower(key) {
					parameters = appendParameter(parameters, key, v)
					continue
				}
			}
		}
	}

	return parameters
}

// extractParameterNames extracts the names of the parameters from the
// cloudformation template body provided. May be either YAML or JSON.
func extractParameterNames(templateBody string) ([]string, error) {
	var (
		data     = []byte(templateBody)
		template struct {
			Parameters map[string]interface{} `yaml:"Parameters"`
		}
	)

	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("unable to unmarshal template body: %w", err)
	}

	var names []string
	for name := range template.Parameters {
		names = append(names, name)
	}
	sort.Strings(names)

	return names, nil
}

func appendParameter(parameters []*cloudformation.Parameter, name, value string) []*cloudformation.Parameter {
	if name == "" {
		panic(fmt.Errorf("illegal attempt to set template parameter with blank name"))
	}
	if value == "" {
		panic(fmt.Errorf("illegal attempt to set template parameter, %v, with blank value", name))
	}

	fmt.Printf("... parameter %v: %v\n", name, value)
	return append(parameters, &cloudformation.Parameter{
		ParameterKey:   aws.String(name),
		ParameterValue: aws.String(value),
	})
}

func uploadLambdas(ctx context.Context, s3API s3iface.S3API, dir string) (map[string]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %w", err)
	}

	keyToID := make(map[string]string)
	for _, fInfo := range fileInfos {
		if strings.HasSuffix(fInfo.Name(), "zip") {
			// ok
			f, err := os.Open(dir + fInfo.Name())
			defer f.Close()
			if err != nil {
				return nil, fmt.Errorf("unable to open file %v: %w", fInfo.Name(), err)
			}
			bucket := config.btLambdasBucket
			key := fInfo.Name()
			fmt.Printf("uploading to s3://%v/%v\n", bucket, key)
			output, err := s3API.PutObjectWithContext(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   f,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to upload file to s3://%v/%v: %w", bucket, key, err)
			}

			if output.VersionId == nil {
				panic(fmt.Errorf("bucket %v must support versioning", bucket))
			}

			prettyKey := cleanKey(key)
			keyToID[prettyKey] = *output.VersionId

		}
	}
	return keyToID, nil
}

// turns something like poller.zip -> pollerzip
func cleanKey(resource string) string {
	var (
		reNotAlphaNumeric = regexp.MustCompile(`[^0-9A-Za-z]+`)
		key               = reNotAlphaNumeric.ReplaceAllString(resource, "")
	)
	key = strings.ToLower(key)
	return key
}
