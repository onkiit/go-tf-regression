package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	log.Println(x)
	log.Println(x)

	pred, err := predict(x, root)
	if err != nil {
		log.Println("predict", err)
		return
	}

	val := tg.Exec(root, []tf.Output{x.Output, y.Output, pred.Output}, nil, &tf.SessionOptions{})

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
		tmp_x, err := strconv.ParseFloat(item.DSPA, 32)
		if err != nil {
			log.Println("strconv")
			return nil, nil, err
		}
		x = append(x, float32(tmp_x))
		y = append(y, item.Displacement)
	}
	xs := tg.NewTensor(root, tg.Const(root, x))
	ys := tg.NewTensor(root, tg.Const(root, y))

	return xs, ys, nil
}

func predict(xs *tg.Tensor, root *op.Scope) (*tg.Tensor, error) {
	var (
		four  = tg.NewTensor(root, tg.Const(root, [1]float32{4}))
		three = tg.NewTensor(root, tg.Const(root, [1]float32{3}))
		a     = tg.NewTensor(root, tg.Const(root, [1]float32{0.00003}))
		b     = tg.NewTensor(root, tg.Const(root, [1]float32{-0.003}))
		c     = tg.NewTensor(root, tg.Const(root, [1]float32{0.11}))
		d     = tg.NewTensor(root, tg.Const(root, [1]float32{-0.9}))
		e     = tg.NewTensor(root, tg.Const(root, [1]float32{4}))
	)
	y := a.Mul(xs.Pow(four.Output).Output).Add(b.Mul(xs.Pow(three.Output).Output).Output).Add(c.Mul(xs.Square().Output).Output).Add(d.Mul(xs.Output).Output).Add(e.Output)
	// y := four
	return y, nil
}
