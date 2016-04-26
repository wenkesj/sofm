package net

import (
	"math"
)

type Node struct {
	Weights []float64
	X       float64
	Y       float64
	Id      int
}

func NewNode(x, y, id int) *Node {
	this := new(Node)
	this.X = float64(x)
	this.Y = float64(y)
	this.Id = id
	return this
}

func (this *Node) GetDistance(inputVector []float64) float64 {
	distance := float64(0)
	for i, weight := range this.Weights {
		distance += math.Pow(inputVector[i]-weight, 2)
	}
	return math.Sqrt(distance)
}

func (this *Node) AdjustWeights(vector []float64, learningRate, influence float64) {
	for i, _ := range this.Weights {
		this.Weights[i] += (vector[i] - this.Weights[i]) * influence * learningRate
	}
}
