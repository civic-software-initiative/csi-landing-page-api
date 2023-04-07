package main

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	from := os.Getenv("SENDER")
	password := os.Getenv("APP_PASSWORD")
	to := []string{os.Getenv("RECEIVER")}

	var message bytes.Buffer
	message.WriteString("Subject: Test Email\r\n")
	message.WriteString("Hello World\r\n")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message.Bytes())
	if err != nil {
		fmt.Println(err)
		return Response{StatusCode: 400}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "",
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "contact-us-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
