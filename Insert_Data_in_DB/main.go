package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Rates struct
type Rates struct {
	Ind int
	INR float64 `json:"INR"`
	CAD float64 `json:"CAD"`
	HKD float64 `json:"HKD"`
}

//Data struct
type Data struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates Rates  `json:"rates"`
}

func hello() {
	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=USD&symbols=INR,CAD,HKD")

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data Data
	json.Unmarshal([]byte(string(body)), &data)
	data.Rates.Ind = 1

	var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))

	rowItem, err := dynamodbattribute.MarshalMap(data.Rates)

	input := &dynamodb.PutItemInput{
		Item:      rowItem,
		TableName: aws.String("guestbook"),
	}

	_, err = db.PutItem(input)

	if err != nil {
		log.Println(err)
		return
	}

	return

}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
	//hello()
}
