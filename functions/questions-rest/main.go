package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/raymonstah/bigtalk/domain/question"
	questionDao "github.com/raymonstah/bigtalk/domain/question/ddb"
	"net/http"
)

type Handler struct {
	dao *questionDao.DAO
}

func (h *Handler) handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case http.MethodGet:
		return h.handleGet(ctx)
	case http.MethodPost:
		return h.handlePost(ctx, request.Body)
	case http.MethodDelete:
		return h.handleDelete(ctx, request.PathParameters[`id`])
	default:
		// not allowed
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       fmt.Sprintf("%v is not allowed", request.HTTPMethod),
		}, fmt.Errorf("%v is not allowed", request.HTTPMethod)
	}

}

func (h *Handler) handleGet(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	questions, err := h.dao.List(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()},
			fmt.Errorf("unable to list questions: %w", err)
	}

	if len(questions) == 0 {
		return emptyList()
	}

	questionsJson, err := json.Marshal(questions)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()},
			fmt.Errorf("unable to marshal message: %w", err)
	}

	return events.APIGatewayProxyResponse{Body: string(questionsJson), StatusCode: http.StatusOK}, nil
}

func (h *Handler) handlePost(ctx context.Context, body string) (events.APIGatewayProxyResponse, error) {

	type createJson struct {
		Question string `json:"question"`
	}

	var inputJson createJson

	err := json.Unmarshal([]byte(body), &inputJson)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "unable to unmarshal request body"},
			fmt.Errorf("unable to unmarshal request body: %w", err)
	}

	q, err := h.dao.Create(ctx, question.CreateInput{
		Question: inputJson.Question,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, fmt.Errorf("unable to create question: %w", err)
	}

	qJson, err := json.Marshal(q)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "unable to marshal json"},
			fmt.Errorf("unable to marshal json: %w", err)
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(qJson)}, nil
}

func (h *Handler) handleDelete(ctx context.Context, s string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func emptyList() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: `[]`}, nil
}

func main() {

	s := session.Must(session.NewSession(aws.NewConfig()))
	questiondao := getQuestionDao(s, "questions")

	handler := Handler{
		dao: questiondao,
	}

	lambda.Start(handler.handle)
}

func getQuestionDao(session *session.Session, questionTableName string) *questionDao.DAO {

	db := dynamo.New(session)
	questionTable := db.Table(questionTableName)
	d := questionDao.New(questionTable)
	return d
}
