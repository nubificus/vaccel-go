// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/nubificus/vaccel-go/vaccel"
)

const DataSize = 30

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Printf("Usage: %s model-dir [iterations]", os.Args[0])
		return
	}
	path := os.Args[1]

	var iters int
	var stat error
	if len(os.Args) == 3 {
		iters, stat = strconv.Atoi(os.Args[2])
		if stat != nil {
			fmt.Println("Error converting argument to int:", stat)
			os.Exit(-1)
		}
	} else {
		iters = 1
	}

	var session vaccel.Session
	var model vaccel.Resource
	var status vaccel.TFStatus
	var inNode vaccel.TFNode
	var outNode vaccel.TFNode
	var runOptions vaccel.TFBuffer
	var inTensor vaccel.TFTensor
	var inTensors []vaccel.TFTensor
	var dataPtr uintptr
	var size uint
	var data []float32

	err := vaccel.ResourceInit(&model, path, vaccel.ResourceModel)
	if err != vaccel.OK {
		fmt.Println("error creating model resource")
		os.Exit(err)
	}

	err = vaccel.SessionInit(&session, 0)
	if err != 0 {
		fmt.Println("error initializing session")
		goto ReleaseResource
	}

	err = vaccel.ResourceRegister(&model, &session)
	if err != 0 {
		fmt.Println("error registering resource with session")
		goto ReleaseSession
	}

	err = vaccel.TFSessionLoad(&session, &model, &status)
	if err != vaccel.OK {
		fmt.Println("Could not load TF session")
		goto UnregisterResource
	}

	err = vaccel.TFStatusRelease(&status)
	if err != vaccel.OK {
		fmt.Println("Could not release TF status")
	}

	err = vaccel.TFNodeInit(&inNode, "serving_default_input_1", 0)
	if err != vaccel.OK {
		fmt.Println("Could not initialize TF Node")
		goto UnregisterResource
	}

	err = vaccel.TFTensorInit(&inTensor, []int64{1, DataSize}, vaccel.TfFloat)
	if err != vaccel.OK {
		fmt.Println("Could not create input tensor")
		goto DeleteTFSession
	}

	data = make([]float32, DataSize)
	for i := 0; i < DataSize; i++ {
		data[i] = 1.0
	}
	dataPtr = uintptr(unsafe.Pointer(&data[0]))
	size = uint(unsafe.Sizeof(data[0]) * DataSize)
	err = vaccel.TFTensorSetData(&inTensor, dataPtr, size, false)
	if err != vaccel.OK {
		fmt.Println("Could not set input tensor data")
		goto DeleteInTensor
	}

	err = vaccel.TFNodeInit(&outNode, "StatefulPartitionedCall", 0)
	if err != vaccel.OK {
		fmt.Println("Cound not configure output TF Node")
		goto DeleteInTensor
	}

	inTensors = []vaccel.TFTensor{inTensor}

	for i := 0; i < iters; i++ {
		outTensors := make([]vaccel.TFTensor, 1)
		err = vaccel.TFSessionRun(
			&session,
			&model,
			&runOptions,
			&inNode,
			inTensors,
			&outNode,
			&outTensors,
			&status,
		)
		if err != vaccel.OK {
			fmt.Println("TF-Session-Run failed")
			goto ReleaseTFStatus
		}

		fmt.Println("Success!")
		fmt.Printf("Output tensor => type:%d nr_dims:%d\n", outTensors[0].Type(),
			outTensors[0].NrDims())

		outDims := outTensors[0].Dims()
		for i := 0; i < outTensors[0].NrDims(); i++ {
			fmt.Printf("dim[%d]: %d\n", i, outDims[i])
		}

		fmt.Println("Result Tensor:")
		outTensors[0].PrintFloat32Data()

		if vaccel.TFTensorRelease(&outTensors[0]) != vaccel.OK {
			fmt.Println("Could not release output tensor")
			goto ReleaseTFStatus
		}

		if i < iters-1 {
			if vaccel.TFStatusRelease(&status) != vaccel.OK {
				fmt.Println("Could not release session run status")
				goto DeleteInTensor
			}
		}
	}

ReleaseTFStatus:
	if vaccel.TFStatusRelease(&status) != vaccel.OK {
		fmt.Println("An error occurred while releasing the status")
	}

DeleteInTensor:
	if vaccel.TFTensorRelease(&inTensor) != vaccel.OK {
		fmt.Println("An error occurred while releasing the tensor")
	}

DeleteTFSession:
	if vaccel.TFSessionDelete(&session, &model, &status) != vaccel.OK {
		fmt.Println("An error occurred while deleting the TF session")
	}

	if vaccel.TFStatusRelease(&status) != vaccel.OK {
		fmt.Println("An error occurred while releasing the status")
	}

UnregisterResource:
	if vaccel.ResourceUnregister(&model, &session) != vaccel.OK {
		fmt.Println("An error occurred while unregistering the resource")
	}

ReleaseSession:
	if vaccel.SessionRelease(&session) != vaccel.OK {
		fmt.Println("An error occurred while releasing the session")
	}

ReleaseResource:
	if vaccel.ResourceRelease(&model) != vaccel.OK {
		fmt.Println("An error occurred while releasing the resource")
	}
}
