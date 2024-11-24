package main

import (
	"github.com/Mohammadmohebi33/toll_calculator/aggregator/client"
	"log"
)

var KafkaTopic = "obudata"
var EndPorint = "http://localhost:8080/aggregate"

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	consumer, err := NewKafkaConsumer(KafkaTopic, svc, client.NewClient(EndPorint))
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start()
}
