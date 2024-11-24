package main

import "log"

var KafkaTopic = "obudata"

func main() {

	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	consumer, err := NewKafkaConsumer(KafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start()
}
