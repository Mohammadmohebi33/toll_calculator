package main

import (
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (receiver *MemoryStore) Insert(distance types.Distance) error {
	receiver.data[distance.OBUID] += distance.Value
	return nil
}

func (receiver *MemoryStore) Get(id int) (float64, error) {
	dist, ok := receiver.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id %d", id)
	}
	return dist, nil
}
