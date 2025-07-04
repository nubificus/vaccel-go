// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <stdlib.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <vaccel.h>

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type TFLiteDataType int

const (
	TfLiteNoType     TFLiteDataType = 0
	TfLiteFloat32    TFLiteDataType = 1
	TfLiteInt32      TFLiteDataType = 2
	TfLiteUint8      TFLiteDataType = 3
	TfLiteInt64      TFLiteDataType = 4
	TfLiteString     TFLiteDataType = 5
	TfLiteBool       TFLiteDataType = 6
	TfLiteInt16      TFLiteDataType = 7
	TfLiteComplex64  TFLiteDataType = 8
	TfLiteInt8       TFLiteDataType = 9
	TfLiteFloat16    TFLiteDataType = 10
	TfLiteFloat64    TFLiteDataType = 11
	TfLiteComplex128 TFLiteDataType = 12
	TfLiteUint64     TFLiteDataType = 13
	TfLiteResource   TFLiteDataType = 14
	TfLiteVariant    TFLiteDataType = 15
	TfLiteUint32     TFLiteDataType = 16
	TfLiteUint16     TFLiteDataType = 17
	TfLiteInt4       TFLiteDataType = 18
)

type TFLiteTensor struct {
	cTFLiteTensor C.struct_vaccel_tflite_tensor
}

func (t *TFLiteTensor) Init(dims []int32, dtype TFLiteDataType) int {
	if t == nil || len(dims) == 0 {
		return EINVAL
	}

	cDims := (*C.int32_t)(C.malloc(C.size_t(len(dims)) * C.size_t(unsafe.Sizeof(C.int32_t(0)))))
	if cDims == nil {
		return ENOMEM
	}
	defer C.free(unsafe.Pointer(cDims))

	base := uintptr(unsafe.Pointer(cDims))
	for i := 0; i < len(dims); i++ {
		ptr := (*C.int32_t)(unsafe.Pointer(base + uintptr(i)*unsafe.Sizeof(C.int32_t(0))))
		*ptr = C.int32_t(dims[i])
	}

	ret := C.vaccel_tflite_tensor_init(
		&t.cTFLiteTensor,
		C.int(len(dims)),
		cDims,
		C.enum_vaccel_tflite_data_type(dtype), //nolint:gocritic
	)

	return int(ret)
}

func (t *TFLiteTensor) Release() int {
	return int(C.vaccel_tflite_tensor_release(&t.cTFLiteTensor)) //nolint:gocritic
}

func (t *TFLiteTensor) Allocate(dims []int32, dtype TFLiteDataType, totalSize uint) int {
	ret := t.Init(dims, dtype)
	if ret != OK {
		return ret
	}

	if totalSize == 0 {
		return OK
	}

	t.cTFLiteTensor.data = C.malloc(C.size_t(totalSize))
	if t.cTFLiteTensor.data == nil {
		C.vaccel_tflite_tensor_release(&t.cTFLiteTensor) //nolint:gocritic
		return ENOMEM
	}

	t.cTFLiteTensor.size = C.size_t(totalSize)
	t.cTFLiteTensor.owned = true

	return OK
}

func (t *TFLiteTensor) SetData(data uintptr, size uint, own bool) int {
	if t == nil {
		return EINVAL
	}

	if t.cTFLiteTensor.data != nil && t.cTFLiteTensor.owned {
		fmt.Println("Previous tensor data will not be freed by release!")
	}

	t.cTFLiteTensor.data = unsafe.Pointer(data)
	t.cTFLiteTensor.size = C.size_t(size)
	t.cTFLiteTensor.owned = C.bool(own)

	return OK
}

func (t *TFLiteTensor) TakeData() (uintptr, uint) {
	if t == nil {
		return 0, 0
	}

	outData := uintptr(t.cTFLiteTensor.data)
	outSize := uint(t.cTFLiteTensor.size)

	t.cTFLiteTensor.data = nil
	t.cTFLiteTensor.size = 0
	t.cTFLiteTensor.owned = false

	return outData, outSize
}

func (t *TFLiteTensor) Data() uintptr {
	if t == nil {
		return 0
	}
	return uintptr(t.cTFLiteTensor.data)
}

func (t *TFLiteTensor) NrDims() int {
	if t == nil || t.cTFLiteTensor.dims == nil || t.cTFLiteTensor.nr_dims <= 0 {
		return -1
	}
	return int(t.cTFLiteTensor.nr_dims)
}

func (t *TFLiteTensor) Type() TFLiteDataType {
	if t == nil || t.cTFLiteTensor.data_type < 1 {
		return -1
	}
	return TFLiteDataType(t.cTFLiteTensor.data_type)
}

func (t *TFLiteTensor) Dims() []int32 {
	if t == nil || t.cTFLiteTensor.dims == nil || t.cTFLiteTensor.nr_dims <= 0 {
		return nil
	}

	n := int(t.cTFLiteTensor.nr_dims)
	ptr := unsafe.Pointer(t.cTFLiteTensor.dims)
	cSlice := unsafe.Slice((*int32)(ptr), n)
	dims := make([]int32, n)
	copy(dims, cSlice)

	return dims
}

func (t *TFLiteTensor) DataPtr() uintptr {
	if t == nil || t.cTFLiteTensor.data == nil {
		return 0
	}
	return uintptr(t.cTFLiteTensor.data)
}

func (t *TFLiteTensor) PrintFloat32Data() {
	if t == nil || t.cTFLiteTensor.data == nil {
		fmt.Println("nil tensor")
		return
	}

	if t.cTFLiteTensor.data_type != C.VACCEL_TFLITE_FLOAT32 {
		fmt.Println("Unsupported data type (only float32 supported)")
		return
	}

	dims := t.Dims()
	if dims == nil {
		fmt.Println("Invalid dims")
		return
	}

	numel := int32(1)
	for _, d := range dims {
		numel *= d
	}

	ptr := (*float32)(unsafe.Pointer(t.DataPtr()))
	slice := unsafe.Slice(ptr, numel)

	fmt.Printf("Tensor shape: %v\n", dims)
	fmt.Println("Values:")
	printRecursiveFloat32TFL(slice, dims, 0)
}

func printRecursiveFloat32TFL(data []float32, dims []int32, level int) {
	if len(dims) == 0 {
		return
	}
	if len(dims) == 1 {
		fmt.Printf("%s[", indent(level))
		for i := int32(0); i < dims[0]; i++ {
			fmt.Printf("%.4f", data[i])
			if i < dims[0]-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")
		return
	}

	size := dims[0]
	sub := int(len(data)) / int(size)
	for i := int32(0); i < size; i++ {
		printRecursiveFloat32TFL(data[i*int32(sub):(i+1)*int32(sub)], dims[1:], level+1)
	}
}

type TFLiteSession struct {
	Sess  *Session
	Model *Resource
}

func (tfls *TFLiteSession) Init(sess *Session, model *Resource) int {
	if sess == nil || model == nil {
		return EINVAL
	}

	tfls.Sess = sess
	tfls.Model = model

	return OK
}

func (tfls *TFLiteSession) Load() int {
	if tfls == nil || tfls.Sess == nil || tfls.Model == nil {
		return EINVAL
	}
	return int(C.vaccel_tflite_session_load(&tfls.Sess.cSess, &tfls.Model.cRes)) //nolint:gocritic
}

func (tfls *TFLiteSession) Run(
	inTensors []TFLiteTensor,
	outTensors *[]TFLiteTensor,
) (int, uint8) {
	if tfls == nil || tfls.Sess == nil || tfls.Model == nil || outTensors == nil {
		return EINVAL, 0
	}

	nrInputs := len(inTensors)
	nrOutputs := len(*outTensors)
	if nrOutputs == 0 {
		return EINVAL, 0
	}

	inBufSize := C.size_t(nrInputs) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cInPtr := C.malloc(inBufSize)
	defer C.free(cInPtr)

	inTensorSlice := unsafe.Slice((**C.struct_vaccel_tflite_tensor)(cInPtr), nrInputs)
	for i := 0; i < nrInputs; i++ {
		inTensorSlice[i] = &inTensors[i].cTFLiteTensor
	}

	outBufSize := C.size_t(nrOutputs) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cOutPtr := C.malloc(outBufSize)
	defer C.free(cOutPtr)

	var cStatus C.uint8_t
	ret := int(C.vaccel_tflite_session_run(
		&tfls.Sess.cSess,
		&tfls.Model.cRes,
		(**C.struct_vaccel_tflite_tensor)(cInPtr),
		C.int(nrInputs),
		(**C.struct_vaccel_tflite_tensor)(cOutPtr),
		C.int(nrOutputs),
		&cStatus, //nolint:gocritic
	))
	if ret != OK {
		return ret, uint8(cStatus)
	}

	outTensorSlice := unsafe.Slice((**C.struct_vaccel_tflite_tensor)(cOutPtr), nrInputs)
	for i := 0; i < nrOutputs; i++ {
		if outTensorSlice[i] != nil {
			nrDims := int(outTensorSlice[i].nr_dims)
			dtype := TFLiteDataType(outTensorSlice[i].data_type)

			dims := make([]int32, nrDims)
			src := unsafe.Slice((*int32)(unsafe.Pointer(outTensorSlice[i].dims)), nrDims)
			copy(dims, src)

			ret = (*outTensors)[i].Init(dims, dtype)
			if ret != OK {
				fmt.Println("Could not initialize output tensor from C dims and type")
				return ret, uint8(cStatus)
			}

			var data unsafe.Pointer
			var size C.size_t
			ret = int(C.vaccel_tflite_tensor_take_data(outTensorSlice[i], &data, &size)) //nolint:gocritic
			if ret != OK {
				fmt.Println("Could not take data from output C tensor")
				return ret, uint8(cStatus)
			}

			ret = int(C.vaccel_tflite_tensor_delete(outTensorSlice[i]))
			if ret != OK {
				fmt.Println("Could not delete output C tensor")
				return ret, uint8(cStatus)
			}

			goData := uintptr(data)
			goSize := uint(size)

			ret = (*outTensors)[i].SetData(goData, goSize, true)
			if ret != OK {
				fmt.Println("Could not set data for Golang output tensor")
				return ret, uint8(cStatus)
			}
		}
	}

	return OK, uint8(cStatus)
}

func (tfls *TFLiteSession) Delete() int {
	err := int(C.vaccel_tflite_session_delete(&tfls.Sess.cSess, &tfls.Model.cRes)) //nolint:gocritic
	if err == OK {
		tfls.Sess = nil
		tfls.Model = nil
	}
	return err
}
