package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

type Data struct {
	Displacement float32 `json:"displacement"`
	DSPA         string  `json:"DSPA"`
}

type Sensor struct {
	Data []Data `json:"data"`
}

func main() {
	// Construct a graph with an operation that produces a string constant.
	s := op.NewScope()
	c := op.Const(s, "Hello from TensorFlow version "+tf.Version())
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}

	xs, ys, err := getTensor()
	if err != nil {
		log.Println("get tensor")
		return
	}
	log.Println(xs)
	log.Println(ys)
	plotData(xs, ys)

	// Execute the graph in a session.
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		panic(err)
	}

	output, err := sess.Run(nil, []tf.Output{c}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(output[0].Value())
}

func getTensor() (*tf.Tensor, *tf.Tensor, error) {
	var (
		dt Sensor
		x  []string
		y  []float32
	)
	sensorFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println("read file")
		return nil, nil, nil
	}

	byteSensor, err := ioutil.ReadAll(sensorFile)
	if err != nil {
		fmt.Println("read all")
		return nil, nil, nil
	}

	if err := json.Unmarshal(byteSensor, &dt); err != nil {
		fmt.Println("unmarshal")
		return nil, nil, nil
	}

	for _, item := range dt.Data {
		x = append(x, item.DSPA)
		y = append(y, item.Displacement)
	}

	xs, err := tf.NewTensor(x)
	if err != nil {
		log.Println("convert x")
		return nil, nil, err
	}

	ys, err := tf.NewTensor(y)
	if err != nil {
		log.Println("convert y")
		return nil, nil, err
	}

	return xs, ys, nil
}

func plotData(xs *tf.Tensor, ys *tf.Tensor) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Polynomial Regression"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	pol := plotter.NewFunction(func(x float64) float64 { return x * x })
	pol.Color = color.RGBA{B: 255, A: 255}

	p.Add(pol)
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.X.Max = 100

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "functions.png"); err != nil {
		panic(err)
	}
}
