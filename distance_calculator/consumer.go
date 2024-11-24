package main

import (
	"encoding/json"
	"github.com/Mohammadmohebi33/toll_calculator/aggregator/client"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.Client) (*KafkaConsumer, error) {
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
		aggClient:   aggClient,
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

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().Unix(),
			OBUID: data.OBUID,
		}

		if err := receiver.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate invoice error %s", err.Error())
			continue
		}

	}
}
