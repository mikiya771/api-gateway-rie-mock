package main

import (
	runtime "github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

type Response struct {
	Message string `json:"msg"`
}

type ResponseType struct {
	Headers    map[string]string
	StatusCode int
	Body       Response
}

func handleRequest() (ResponseType, error) {
	return ResponseType{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: http.StatusOK,
		Body: Response{
			Message: "Hello World",
		},
	}, nil
}

func main() {
	runtime.Start(handleRequest)
}
