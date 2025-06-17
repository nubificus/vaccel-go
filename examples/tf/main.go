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
	var tfSess vaccel.TFSession
	var status vaccel.TFStatus
	var inNode vaccel.TFNode
	var outNode vaccel.TFNode
	var runOptions vaccel.TFBuffer
	var inTensor vaccel.TFTensor
	var inTensors []vaccel.TFTensor
	var dataPtr uintptr
	var size uint
	var err int
	var data []float32

	err = model.Init(path, vaccel.ResourceModel)
	if err != vaccel.OK {
		fmt.Println("error creating model resource")
		os.Exit(err)
	}

	err = session.Init(0)
	if err != 0 {
		fmt.Println("error initializing session")
		goto ReleaseResource
	}

	err = session.Register(&model)
	if err != 0 {
		fmt.Println("error registering resource with session")
		goto ReleaseSession
	}

	tfSess = vaccel.TFSession{Sess: &session, Model: &model}
	err = tfSess.Load(&status)
	if err != vaccel.OK {
		fmt.Println("Could not load TF session")
		goto UnregisterResource
	}

	err = status.Release()
	if err != vaccel.OK {
		fmt.Println("Could not release TF status")
	}

	err = inNode.Init("serving_default_input_1", 0)
	if err != vaccel.OK {
		fmt.Println("Could not initialize TF Node")
		goto UnregisterResource
	}

	err = inTensor.Init([]int64{1, DataSize}, vaccel.TfFloat)
	if err != vaccel.OK {
		fmt.Println("Could not create input tensor")
		goto DeleteTFSession
	}

	data = make([]float32, DataSize)
	for i := 0; i < DataSize; i++ {
		data[i] = float32(iters)
	}
	dataPtr = uintptr(unsafe.Pointer(&data[0]))
	size = uint(unsafe.Sizeof(data[0]) * DataSize)
	err = inTensor.SetData(dataPtr, size, false)
	if err != vaccel.OK {
		fmt.Println("Could not set input tensor data")
		goto DeleteInTensor
	}

	err = outNode.Init("StatefulPartitionedCall", 0)
	if err != vaccel.OK {
		fmt.Println("Cound not configure output TF Node")
		goto DeleteInTensor
	}

	inTensors = []vaccel.TFTensor{inTensor}

	for i := 0; i < iters; i++ {
		outTensors := make([]vaccel.TFTensor, 1)
		err = tfSess.Run(
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

		if outTensors[0].Release() != vaccel.OK {
			fmt.Println("Could not release output tensor")
			goto ReleaseTFStatus
		}

		if i < iters-1 {
			if status.Release() != vaccel.OK {
				fmt.Println("Could not release session run status")
				goto DeleteInTensor
			}
		}
	}

	if inNode.Release() != vaccel.OK {
		fmt.Println("Could not release input node")
	}

	if outNode.Release() != vaccel.OK {
		fmt.Println("Could not release output node")
	}

ReleaseTFStatus:
	if status.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the status")
	}

DeleteInTensor:
	if inTensor.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the tensor")
	}

DeleteTFSession:
	if tfSess.Delete(&status) != vaccel.OK {
		fmt.Println("An error occurred while deleting the TF session")
	}

	if status.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the status")
	}

UnregisterResource:
	if session.Unregister(&model) != vaccel.OK {
		fmt.Println("An error occurred while unregistering the resource")
	}

ReleaseSession:
	if session.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the session")
	}

ReleaseResource:
	if model.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the resource")
	}
}
