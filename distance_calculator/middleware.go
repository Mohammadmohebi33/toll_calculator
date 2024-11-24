package main

import (
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{next: next}
}

func (l *LogMiddleware) CalculatorDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info("calculator distance")
	}(time.Now())
	dist, err = l.next.CalculatorDistance(data)
	return
}
