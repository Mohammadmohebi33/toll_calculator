package main

import (
	"github.com/Mohammadmohebi33/toll_calculator/aggregator/client"
	"log"
)

var KafkaTopic = "obudata"
var HttpEndPoint = "http://localhost:8080/aggregate"
var GrpcEndPoint = ":8081"

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	grpcClient, err := client.NewGRPCClient(GrpcEndPoint)
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}

	consumer, err := NewKafkaConsumer(KafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start()
}
