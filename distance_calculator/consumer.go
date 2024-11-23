package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
	}, nil
}

func (receiver *KafkaConsumer) Start() {
	logrus.Info("kafka consumer started")
	receiver.isRunning = true
	receiver.readMessageLoop()
}

func (receiver *KafkaConsumer) readMessageLoop() {
	for receiver.isRunning {
		msg, err := receiver.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error")
			continue
		}
		var data types.OBUData

		if err = json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("json unmarshal error %s", err.Error())
			continue
		}

		distance, err := receiver.calcService.CalculatorDistance(data)
		if err != nil {
			logrus.Errorf("calc service error %s", err.Error())
		}
		fmt.Printf("distance %.2f\n", distance)
	}
}
