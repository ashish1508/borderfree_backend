package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ConnectionItem struct {
	ConnectionID string `json:"connectionID"`
}

const APIGatewayEndpoint = "https://v2cihxqx2i.execute-api.ap-south-1.amazonaws.com/dev/"

func NewDynamoDBSession() *dynamodb.DynamoDB {
	return dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))
}

func NewAPIGatewaySession() *apigatewaymanagementapi.ApiGatewayManagementApi {
	sess, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String(APIGatewayEndpoint),
	})
	return apigatewaymanagementapi.New(sess)
}
