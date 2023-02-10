package main

import (
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"
	api "github.com/chenliu1993/ec2simplelambda/pkg"
)

func main() {
	log.Printf("begin")
	runtime.Start(api.HandlerReq)
	log.Printf("end")
}
