// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/nfnt/resize"
	"github.com/nubificus/vaccel-go/vaccel"
)

const (
	ImageWidth    = 224
	ImageHeight   = 224
	ImageChannels = 3
)

func loadLabels(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var labels []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			labels = append(labels, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return labels, nil
}

func preprocessImage(img image.Image) ([]float32, error) {
	resized := resize.Resize(ImageWidth, ImageHeight, img, resize.NearestNeighbor)
	if resized.Bounds().Dx() != ImageWidth || resized.Bounds().Dy() != ImageHeight {
		return nil, errors.New("resized image has wrong dimensions")
	}

	mean := [3]float32{0.485, 0.456, 0.406}
	std := [3]float32{0.229, 0.224, 0.225}

	tensor := make([]float32, ImageChannels*ImageHeight*ImageWidth)
	for c := 0; c < ImageChannels; c++ {
		for y := 0; y < ImageHeight; y++ {
			for x := 0; x < ImageWidth; x++ {
				r, g, b, _ := resized.At(x, y).RGBA()
				pixel := []float32{
					float32(r>>8) / 255.0,
					float32(g>>8) / 255.0,
					float32(b>>8) / 255.0,
				}
				idx := c*ImageHeight*ImageWidth + y*ImageWidth + x
				tensor[idx] = (pixel[c] - mean[c]) / std[c]
			}
		}
	}
	return tensor, nil
}

func loadAndPreprocessImage(path string) ([]float32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}

	return preprocessImage(img)
}

func processResult(output []float32, labelFile string) error {
	if len(output) == 0 {
		return errors.New("empty output")
	}
	labels, err := loadLabels(labelFile)
	if err != nil {
		return err
	}
	if len(labels) > len(output) {
		return errors.New("output smaller than label set")
	}

	maxIdx := 0
	maxVal := output[0]
	for i := 1; i < len(labels); i++ {
		if output[i] > maxVal {
			maxVal = output[i]
			maxIdx = i
		}
	}

	fmt.Printf("Prediction: %s\n", labels[maxIdx])
	return nil
}

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Printf("Usage: %s <image_file> <model_file> <labels_file> [iterations]\n", os.Args[0])
		os.Exit(1)
	}

	imageFile := os.Args[1]
	modelFile := os.Args[2]
	labelsFile := os.Args[3]
	iterations := 1
	if len(os.Args) == 5 {
		if val, err := strconv.Atoi(os.Args[4]); err == nil {
			iterations = val
		}
	}

	var session vaccel.Session
	var model vaccel.Resource
	var runOptions vaccel.TorchBuffer
	var inTensor vaccel.TorchTensor
	var inTensors []vaccel.TorchTensor
	var inputSize int
	var inputData []float32
	var stat error

	err := model.Init(modelFile, vaccel.ResourceModel)
	if err != vaccel.OK {
		fmt.Println("error initializing resource")
		os.Exit(err)
	}

	err = session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		goto ReleaseResource
	}

	err = session.Register(&model)
	if err != vaccel.OK {
		fmt.Println("Could not register model")
		goto ReleaseSession
	}

	err = inTensor.Init([]int64{1, ImageChannels, ImageWidth, ImageHeight}, vaccel.TorchFloat)
	if err != vaccel.OK {
		fmt.Println("Could not create input tensor")
		goto UnregisterResource
	}

	inputSize = int(unsafe.Sizeof(float32(0))) * int(ImageChannels*ImageWidth*ImageHeight)
	if inputSize < 0 {
		fmt.Println("inputSize must be non-negative")
		goto ReleaseInTensor
	}

	inputData, stat = loadAndPreprocessImage(imageFile)
	if stat != nil {
		fmt.Fprintf(os.Stderr, "Could not load and preprocess image: %v\n", err)
		goto ReleaseInTensor
	}

	err = inTensor.SetData(uintptr(unsafe.Pointer(&inputData[0])), uint(inputSize), false)
	if err != vaccel.OK {
		fmt.Println("Could not set tensor data")
		goto ReleaseInTensor
	}

	inTensors = []vaccel.TorchTensor{inTensor}

	err = vaccel.TorchModelLoad(&session, &model)
	if err != vaccel.OK {
		fmt.Println("Could not load model")
		goto ReleaseInTensor
	}

	for i := 0; i < iterations; i++ {
		outTensors := make([]vaccel.TorchTensor, 1)

		err = vaccel.TorchModelRun(&session, &model, &runOptions, inTensors, &outTensors)
		if err != vaccel.OK {
			fmt.Println("TorchModelRun failed")
			goto ReleaseInTensor
		}

		outData, outLen := outTensors[0].TakeData()
		if outData == 0 || outLen == 0 {
			fmt.Println("Could not take data from out tensor")
			break
		}

		output := make([]float32, 1000)
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&output))
		sliceHeader.Data = outData
		sliceHeader.Len = 1000
		sliceHeader.Cap = 1000

		fmt.Println("Success!")

		stat := processResult(output, labelsFile)
		if stat != nil {
			fmt.Println("Could not process result")
			break
		}

		err = outTensors[0].Release()
		if err != vaccel.OK {
			fmt.Println("Could not release out tensor")
			break
		}
	}

ReleaseInTensor:
	err = inTensor.Release()
	if err != vaccel.OK {
		fmt.Println("Could not release tensor")
	}
UnregisterResource:
	err = session.Unregister(&model)
	if err != vaccel.OK {
		fmt.Println("Could not unregister model")
	}
ReleaseSession:
	err = session.Release()
	if err != vaccel.OK {
		fmt.Println("Could not release session")
	}
ReleaseResource:
	err = model.Release()
	if err != vaccel.OK {
		fmt.Println("Could not release model")
		os.Exit(err)
	}
}
