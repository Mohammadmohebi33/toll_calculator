package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var kafkaTopic = "obudata"

func main() {

	recv, _ := NewDataReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":8080", nil)

}

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
	prod *kafka.Producer
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &DataReceiver{
		msg:  make(chan types.OBUData, 128),
		prod: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Produce messages to topic (asynchronously)
	topic := kafkaTopic
	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

	return err
}

func (dr DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	dr.conn = conn

	go dr.wsReceiverLoop()
}

func (dr DataReceiver) wsReceiverLoop() {
	fmt.Println("obu connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error", err)
			continue
		}

		fmt.Println("data receive :", data)
		if err := dr.produceData(data); err != nil {
			log.Println("produce error", err)
		}
	}
}
