package model

import (
	"GoSlicer/util"
	"errors"
	"io"
	"os"
	"strings"
)

type SimpleModel struct {
}

func loadModel(filename string, matrix util.FMatrix3x3) (*SimpleModel, error) {
	splitted := strings.Split(filename, ".")
	if len(splitted) <= 1 {
		return nil, errors.New("the file has no extension")
	}

	extension := splitted[len(splitted)-1]

	if extension == "stl" {
		return loadModelSTL(filename, matrix)
	}
	return nil, errors.New("the file is not a stl file")
}

func loadModelSTL(filename string, matrix util.FMatrix3x3) (*SimpleModel, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	var header [5]byte
	_, err = io.ReadFull(r, header[:])
	if err != nil {
		return nil, err
	}

	solidHeader := []byte("SOLID")
	for i, b := range header {
		if b != solidHeader[i] {
			return loadModelSTLBinary(filename, matrix)
		}
	}

	return loadModelSTLAscii(filename, matrix)
}

func loadModelSTLAscii(filename string, matrix util.FMatrix3x3) (*SimpleModel, error) {
	m := SimpleModel{}
	r, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error while opening file")
	}

	defer r.Close()

	var data []byte
	io.ReadFull(r, data)

	return nil, nil
}

func loadModelSTLBinary(filename string, matrix util.FMatrix3x3) (*SimpleModel, error) {
	return nil, errors.New("not implemented")
}
