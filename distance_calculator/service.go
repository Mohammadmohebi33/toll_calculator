package main

import (
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"math"
)

type CalculatorServicer interface {
	CalculatorDistance(data types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		prevPoint: make([]float64, 0),
	}
}

func (s *CalculatorService) CalculatorDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoint) > 0 {
		distance = calculatorDistande(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculatorDistande(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
