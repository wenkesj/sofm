package net

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"
)

type Match struct {
	Value float64
	Node  *Node
}

type Network struct {
	MapRadius float64
	Nodes     []*Node
	Size      int
}

func NewNetwork(w, h int) *Network {
	this := new(Network)
	this.Size = w * h
	x := 0
	this.Nodes = make([]*Node, this.Size)
	for i := 0; i < this.Size; i++ {
		y := i % w
		if y == 0 && i != 0 {
			x++
		}
		this.Nodes[i] = NewNode(x, y, i)
	}
	this.MapRadius = math.Max(float64(w), float64(h)) / 2
	return this
}

func Save(network *Network, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	dataEncoder := gob.NewEncoder(file)
	dataEncoder.Encode(network)
	return nil
}

func Load(path string) (*Network, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	network := new(Network)
	dataDecoder := gob.NewDecoder(file)
	err = dataDecoder.Decode(network)

	if err != nil {
		return nil, err
	}

	return network, nil
}

func (this *Network) Probe(vector []float64) *Node {
	bestMatchingNode := this.Nodes[0]
	min := bestMatchingNode.GetDistance(vector)

	for _, node := range this.Nodes {
		distance := node.GetDistance(vector)
		if distance <= min {
			min = distance
			bestMatchingNode = node
		}
	}

	return bestMatchingNode
}

func (this *Network) Threshold() {
	for _, node := range this.Nodes {
		adjustedWeights := make([]float64, len(node.Weights))
		for j, weight := range node.Weights {
			if weight > 0 {
				adjustedWeights[j] = WeightMax
			} else {
				adjustedWeights[j] = WeightMin
			}
		}
		node.Weights = adjustedWeights
	}
}

func (this *Network) InitializeWeights(data [][]float64) {
	precision := math.Pow(10, (math.Ceil(math.Log(float64(this.Size))/math.Ln10) + 2))
	scale := 1

	for _, node := range this.Nodes {
		node.Weights = make([]float64, len(data[0]))
		for j, _ := range data[0] {
			node.Weights[j] = Round(rand.Float64()*precision, 2) / precision * float64(scale)
		}
	}
}

func (this *Network) Train(data [][]float64, iterations int, learningRate float64) *Network {
	var vector []float64
	recycle := make([][]float64, len(data))
	copy(recycle, data)

	this.InitializeWeights(data)

	timeConstant := float64(iterations) / math.Log(this.MapRadius)
	constantLearningRate := learningRate

	recycle = make([][]float64, len(data))
	copy(recycle, data)

	for iteration := 0; iteration < iterations; iteration++ {
		if len(recycle) <= 1 {
			recycle = make([][]float64, len(data))
			copy(recycle, data)
		}

		vector, recycle = PluckVector(recycle)

		winningNode := this.Probe(vector)
		neighbourhoodRadius := this.MapRadius * math.Exp(-float64(iteration)/timeConstant)

		for _, node := range this.Nodes {
			squaredDistance := math.Pow(winningNode.X-node.X, 2) + math.Pow(winningNode.Y-node.Y, 2)
			squaredRadius := math.Pow(neighbourhoodRadius, 2)

			if squaredDistance < squaredRadius {
				node.AdjustWeights(vector, learningRate, math.Exp(-squaredDistance/(2*squaredRadius)))
			}
		}
		learningRate = constantLearningRate * math.Exp(-float64(iteration)/float64(iterations))
	}
	return this
}
