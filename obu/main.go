package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var sendInterval = time.Second

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

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
	obuIDs := generateOBUSIDS(20)
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLocation()
			data := OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Println(data)
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
