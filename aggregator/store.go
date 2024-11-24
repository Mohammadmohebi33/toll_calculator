package main

import "github.com/Mohammadmohebi33/toll_calculator/types"

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
