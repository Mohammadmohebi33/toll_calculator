package main

import "log"

var KafkaTopic = "obudata"

func main() {
	consumer, err := NewKafkaConsumer(KafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	consumer.Start()
}
