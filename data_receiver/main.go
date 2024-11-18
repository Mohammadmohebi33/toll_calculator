package main

import (
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":8080", nil)

}

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msg: make(chan types.OBUData, 128),
	}
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

		fmt.Printf("received obu date from [%d] :: <lat %2f, long %2f >\n", data.OBUID, data.Lat, data.Long)
		dr.msg <- data
	}
}
