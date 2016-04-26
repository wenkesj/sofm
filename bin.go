package main

import (
	"bufio"
	"bytes"
	"github.com/wenkesj/sofm/net"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	app            = kingpin.New("sofm", "Self organized feature mapping tool")
	width          = app.Flag("width", "Width of the grid").Required().Int()
	height         = app.Flag("height", "Height of the grid").Required().Int()
	iterations     = app.Flag("iterations", "Training iterations").Required().Int()
	learningRate   = app.Flag("learning-rate", "Learning rate").Required().Float()
	dataPath       = app.Flag("train-data", "Data file").Required().String()
	labelsPath     = app.Flag("train-labels", "Label file").Required().String()
	outPath        = app.Flag("output", "Output file").Default("").String()
	save           = app.Flag("save", "Save network file").Default("").String()
	load           = app.Flag("load", "Load network file").Default("").String()
	testDataPath   = app.Flag("test-data", "Load network file").Default("").String()
	testLabelsPath = app.Flag("test-labels", "Load network file").Default("").String()
	threshold      = app.Flag("threshold", "Normalize final node weights").Default("false").Bool()
	normalize      = app.Flag("normalize", "Normalize inputs").Default("false").Bool()
)

func ReadLabels(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func ReadLines(path string) (lines [][]float64, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			line := strings.Fields(buffer.String())
			intLine := make([]float64, len(line))
			for i, val := range line {
				intLine[i], err = strconv.ParseFloat(val, 64)
				if err != nil {
					panic(err)
				}
			}
			lines = append(lines, intLine)
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func RunProbe(testLabels []string, testData [][]float64, network *net.Network, buffer *bytes.Buffer) {
	keys := make([]string, len(network.Nodes))
	for i, _ := range keys {
		keys[i] = "-"
	}
	for i, key := range testLabels {
		vector := testData[i]
		minId := network.Nodes[0].Id
		minDistance := network.Nodes[0].GetDistance(vector)
		for _, node := range network.Nodes {
			distance := node.GetDistance(vector)
			if distance <= minDistance {
				minDistance = distance
				minId = node.Id
			}
		}
		keys[minId] = key
	}
	for j, _ := range network.Nodes {
		key := keys[j]
		if key != "-" {
			if (j+1)%*width == 0 {
				buffer.WriteString(key + "\n")
			} else {
				buffer.WriteString(key + "\t")
			}
		} else {
			if (j+1)%*width == 0 {
				buffer.WriteString(key + "\n")
			} else {
				buffer.WriteString(key + "\t")
			}
		}
	}
}

func RunMatrix(labels []string, data [][]float64, network *net.Network, buffer *bytes.Buffer) {
	labelMap := make(map[string][]float64)
	for i, key := range labels {
		labelMap[key] = data[i]
	}
	var minKey string
	var minDistance float64
	for j, node := range network.Nodes {
		i := 0
		for key, value := range labelMap {
			if i == 0 {
				minKey = key
				minDistance = node.GetDistance(value)
				i++
			}
			distance := node.GetDistance(value)
			if distance <= minDistance {
				minKey = key
				minDistance = distance
			}
		}
		if (j+1)%*width == 0 {
			buffer.WriteString(minKey + "\n")
		} else {
			buffer.WriteString(minKey + "\t")
		}
	}
}

func Normalize(value float64) float64 {
	if value > 0 {
		return net.WeightMax
	} else {
		return net.WeightMin
	}
}

func Scale(data [][]float64) [][]float64 {
	scaled := make([][]float64, len(data))
	for i, vector := range data {
		scaled[i] = make([]float64, len(vector))
		for j, value := range vector {
			scaled[i][j] = Normalize(value)
		}
	}
	return scaled
}

func main() {
	app.Parse(os.Args[1:])
	data, err := ReadLines(*dataPath)
	labels, err := ReadLabels(*labelsPath)
	if err != nil {
		panic(err)
	}

	if *normalize {
		data = Scale(data)
	}

	var network *net.Network
	if *load != "" {
		network, err = net.Load(*load)
		if err != nil {
			panic(err)
		}
	} else {
		network = net.NewNetwork(*width, *height)
	}

	network.Train(data, *iterations, *learningRate)

	if *threshold {
		network.Threshold()
	}

	if *save != "" {
		err = net.Save(network, *save)
		if err != nil {
			panic(err)
		}
	}

	if *outPath != "" {
		buffer := new(bytes.Buffer)
		buffer.WriteString("Train Response Matrix\n")
		RunProbe(labels, data, network, buffer)
		buffer.WriteString("\n\n")
		buffer.WriteString("Train Response Matrix\n")
		RunMatrix(labels, data, network, buffer)
		buffer.WriteString("\n\n")

		if *testDataPath != "" && *testLabelsPath != "" {
			testData, err := ReadLines(*testDataPath)
			testLabels, err := ReadLabels(*testLabelsPath)
			if err != nil {
				panic(err)
			}
			testData = Scale(testData)
			buffer.WriteString("Test Response Matrix\n")
			RunProbe(testLabels, testData, network, buffer)
			buffer.WriteString("\n\n")
			buffer.WriteString("Test Response Matrix\n")
			RunMatrix(testLabels, testData, network, buffer)
		}

		file, err := os.OpenFile(*outPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		file.WriteString(buffer.String())
		file.Close()
	}
}
