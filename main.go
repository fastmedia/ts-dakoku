package main

import (
	"github.com/ngs/ts-dakoku/app"
	"github.com/aws/aws-lambda-go/lambda"
)
func main() {
	server, err := app.Run()
	if err != nil {
		panic(err)
	}
	lambda.Start(server)
}
