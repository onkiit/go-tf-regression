package main

import (
	"encoding/json"
	"fmt"
	"go-tf-regression/helper"
	"log"
	"net/http"
	"strconv"

	tg "github.com/galeone/tfgo"
	"github.com/gin-gonic/gin"
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

type OriginalResponse struct {
	XS   []float32 `json:"xs,omitempty"`
	YS   []float32 `json:"ys,omitempty"`
	Pred []float32 `json:"predict,omitempty"`
}

const (
	iteration = 100
)

func main() {
	/*
	*get data sensor v
	*convert to tensor v
	*create variable coefficient
	*create predict function using polynomial equation
	*create lost function
	*create optimizer using Adam
	*training data
	*send curv data using REST API and plot to python
	*send prediction of the coefficients
	 */
	router := gin.Default()
	root := tg.NewRoot()
	x, y, err := getTensor(root)
	if err != nil {
		log.Println("get tensor")
		return
	}

	pred, err := predict(x, root)
	if err != nil {
		log.Println("predict", err)
		return
	}

	val := tg.Exec(root.SubScope("execution"), []tf.Output{x.Output, y.Output, pred.Output}, nil, &tf.SessionOptions{})

	xs := val[0].Value()
	ys := val[1].Value()
	ps := val[2].Value()

	router.GET("/original", func(c *gin.Context) {
		respData := &OriginalResponse{
			XS:   (xs).([]float32),
			YS:   (ys).([]float32),
			Pred: (ps).([]float32),
		}
		data, err := json.Marshal(respData)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, string(data))
	})

	router.Run(":8001")

}

func getTensor(root *op.Scope) (*tg.Tensor, *tg.Tensor, error) {
	var (
		dt Sensor
		x  []float32
		y  []float32
	)
	filename := "data.json"

	byteSensor, err := helper.ReadData(filename)
	if err != nil {
		log.Println("byte err")
		return nil, nil, err
	}

	if err := json.Unmarshal(byteSensor, &dt); err != nil {
		fmt.Println("unmarshal")
		return nil, nil, nil
	}

	for _, item := range dt.Data {
		tmp_x, err := strconv.ParseFloat(item.DSPA, 32)
		if err != nil {
			log.Println("strconv")
			return nil, nil, err
		}
		x = append(x, float32(tmp_x))
		y = append(y, item.Displacement)
	}
	xs := tg.NewTensor(root.SubScope("xs"), tg.Const(root.SubScope("x"), x))
	ys := tg.NewTensor(root.SubScope("ys"), tg.Const(root.SubScope("y"), y))

	return xs, ys, nil
}

func predict(xs *tg.Tensor, root *op.Scope) (*tg.Tensor, error) {
	var (
		four  = tg.NewTensor(root.SubScope("pow"), tg.Const(root.SubScope("const"), [1]float32{4}))
		three = tg.NewTensor(root.SubScope("pow"), tg.Const(root.SubScope("const"), [1]float32{3}))
		a     = tg.NewTensor(root.SubScope("coef"), tg.Const(root.SubScope("const"), [1]float32{0.00003}))
		b     = tg.NewTensor(root.SubScope("coef"), tg.Const(root.SubScope("const"), [1]float32{-0.003}))
		c     = tg.NewTensor(root.SubScope("coef"), tg.Const(root.SubScope("const"), [1]float32{0.11}))
		d     = tg.NewTensor(root.SubScope("coef"), tg.Const(root.SubScope("const"), [1]float32{-0.9}))
		e     = tg.NewTensor(root.SubScope("coef"), tg.Const(root.SubScope("const"), [1]float32{4}))
	)
	y := a.Mul(xs.Pow(four.Output).Output).Add(b.Mul(xs.Pow(three.Output).Output).Output).Add(c.Mul(xs.Square().Output).Output).Add(d.Mul(xs.Output).Output).Add(e.Output)

	return y, nil
}

func getLoss(predict *tg.Tensor, label *tg.Tensor) *tg.Tensor {
	loss := predict.Substract(label.Output)
	return loss
}

func trainData(xs *tg.Tensor, ys *tg.Tensor, root *op.Scope) {
	learningRate := op.Const(root.SubScope("rate"), [1]float32{0.01})
	for i := 0; i < iteration; i++ {
		pred, err := predict(xs, root)
		if err != nil {
			log.Println("predict", err)
			return
		}

		_ = op.ResourceApplyAdam(root.SubScope("optimizer"), pred.Output, pred.Output, pred.Output, pred.Output, pred.Output, learningRate, pred.Output, pred.Output, pred.Output, pred.Output)
	}
}
