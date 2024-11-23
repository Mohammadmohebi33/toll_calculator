package main

import (
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	recv, _ := NewDataReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":8081", nil)
}

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
	prod DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)

	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		msg:  make(chan types.OBUData, 128),
		prod: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
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

		if err := dr.produceData(data); err != nil {
			log.Println("produce error", err)
		}
	}
}
