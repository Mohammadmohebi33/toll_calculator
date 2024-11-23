package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
}

func NewKafkaConsumer(topic string) (*KafkaConsumer, error) {
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
		consumer: c,
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
		fmt.Println(msg)
	}
}
