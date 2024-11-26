package main

import (
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/sirupsen/logrus"
)

const bestPrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(id int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (a *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuid":    distance.OBUID,
		"distance": distance.Value,
		"unix":     distance.Unix,
	}).Info("aggregate invoice")
	return a.store.Insert(distance)
}

func (a *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := a.store.Get(obuID)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   bestPrice * dist,
	}

	return inv, nil
}
