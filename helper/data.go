package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	x []string
	y []float32
)

func ReadData(filename string, dt interface{}) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("read file")
		return nil, err
	}

	byteFile, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read all")
		return nil, err
	}

	if err := json.Unmarshal(byteFile, &dt); err != nil {
		fmt.Println("unmarshal")
		return nil, err
	}

	return dt, nil
}
