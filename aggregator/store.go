package main

import "github.com/Mohammadmohebi33/toll_calculator/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (receiver *MemoryStore) Insert(distance types.Distance) error {
	return nil
}
