package helper

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadData(filename string) ([]byte, error) {
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

	return byteFile, nil
}
