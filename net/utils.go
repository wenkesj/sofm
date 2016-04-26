package net

import (
	"math"
	"math/rand"
	"time"
)

var (
	WeightMin = -1.0
	WeightMax = 1.0
)

func RandomRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func PluckVector(vectors [][]float64) ([]float64, [][]float64) {
	randomNumber := RandomRange(0, len(vectors)-1)
	plucked := vectors[randomNumber]
	vectors = append(vectors[:randomNumber], vectors[randomNumber+1:]...)
	return plucked, vectors
}

func RandomVector(vectors [][]float64) []float64 {
	randomNumber := RandomRange(0, len(vectors)-1)
	return vectors[randomNumber]
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	intermed += .5
	x = .5
	if frac < 0.0 {
		x = -.5
		intermed -= 1
	}
	if frac >= x {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
