package main

import (
	"context"
	"github.com/Mohammadmohebi33/toll_calculator/aggregator/client"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"log"
	"time"
)

func main() {
	c, err := client.NewGRPCClient(":8081")

	if err != nil {
		log.Fatal(err)
	}

	if _, err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 60.1,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}

}
