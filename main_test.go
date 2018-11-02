package main

import (
	"reflect"
	"testing"

	tg "github.com/galeone/tfgo"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func Test_getTensor(t *testing.T) {
	type args struct {
		root *op.Scope
	}
	tests := []struct {
		name    string
		args    args
		want    *tg.Tensor
		want1   *tg.Tensor
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getTensor(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTensor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTensor() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getTensor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_predict(t *testing.T) {
	type args struct {
		xs   *tg.Tensor
		root *op.Scope
	}
	tests := []struct {
		name    string
		args    args
		want    *tg.Tensor
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := predict(tt.args.xs, tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("predict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
