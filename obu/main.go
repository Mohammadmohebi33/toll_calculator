package main

import (
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"time"
)

const wsEndpoint = "ws://127.0.0.1:8081/ws"

var sendInterval = time.Second * 60

func genLocation() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return f + n
}

func generateOBUSIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	obuIDs := generateOBUSIDS(1)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLocation()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Println("write:", err)
			}
			fmt.Println(data.OBUID, data.Lat, data.Long)
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
