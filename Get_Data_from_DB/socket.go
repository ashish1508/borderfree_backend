package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	lambda.Start(HandleConnect)
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	dynamodbSession := NewDynamoDBSession()
	connectionItem := ConnectionItem{
		ConnectionID: request.RequestContext.ConnectionID,
	}
	attributeValues, _ := dynamodbattribute.MarshalMap(connectionItem)

	//log.Println(request.RequestContext.EventType)

	if request.RequestContext.EventType == "CONNECT" {

		input := &dynamodb.PutItemInput{
			Item:      attributeValues,
			TableName: aws.String("connections"),
		}
		dynamodbSession.PutItem(input)

	} else if request.RequestContext.EventType == "DISCONNECT" {

		input := &dynamodb.DeleteItemInput{
			Key:       attributeValues,
			TableName: aws.String("connections"),
		}
		dynamodbSession.DeleteItem(input)

	} else {

		inputconn := &dynamodb.ScanInput{
			TableName: aws.String("connections"),
		}
		conn, _ := dynamodbSession.Scan(inputconn)

		inputrates := &dynamodb.ScanInput{
			TableName: aws.String("guestbook"),
		}
		rates, _ := dynamodbSession.Scan(inputrates)

		var output []ConnectionItem
		dynamodbattribute.UnmarshalListOfMaps(conn.Items, &output)

		apigatewaySession := NewAPIGatewaySession()

		jsonData, _ := json.Marshal(rates.Items)

		for _, item := range output {
			connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(item.ConnectionID),
				Data:         jsonData,
			}
			_, err := apigatewaySession.PostToConnection(connectionInput)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
