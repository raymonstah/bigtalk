package main

import (
	"context"
	"encoding/json"
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

func (h *Handler) handle(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

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
		}
	}

	h.dao.Get()
	json.Marshal()
	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Body:              response,
		IsBase64Encoded:   false,
	}

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
