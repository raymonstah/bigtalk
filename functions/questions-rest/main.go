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
	questionDao "github.com/raymonstah/bigtalk/domain/question/ddb"
	"net/http"
)

type Handler struct {
	dao *questionDao.DAO
}

func (h *Handler) handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	case http.MethodDelete:
		break
	default:
		// not allowed
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       fmt.Sprintf("%v is not allowed", request.HTTPMethod),
		}, fmt.Errorf("%v is not allowed", request.HTTPMethod)
	}

	questions, err := h.dao.List(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()},
			fmt.Errorf("unable to list questions: %w", err)
	}

	questionsJson, err := json.Marshal(questions)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()},
			fmt.Errorf("unable to marshal message: %w", err)
	}

	return events.APIGatewayProxyResponse{Body: string(questionsJson), StatusCode: http.StatusOK}, nil

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
