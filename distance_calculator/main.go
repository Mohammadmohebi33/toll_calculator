package main

import (
	"log"

	"github.com/Mohammadmohebi33/toll_calculator/aggregator/client"
)

var KafkaTopic = "obudata"
var HttpEndPoint = "http://localhost:3000"

//var GrpcEndPoint = ":8081"

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// grpcClient, err := client.NewGRPCClient(GrpcEndPoint)
	// if err != nil {
	// 	log.Fatalf("failed to create gRPC client: %v", err)
	// }

	consumer, err := NewKafkaConsumer(KafkaTopic, svc, client.NewHttpClient(HttpEndPoint))
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start()
}
