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

type TorchDataType int

const (
	TorchByte  TorchDataType = 1
	TorchChar  TorchDataType = 2
	TorchShort TorchDataType = 3
	TorchInt   TorchDataType = 4
	TorchLong  TorchDataType = 5
	TorchHalf  TorchDataType = 6
	TorchFloat TorchDataType = 7
)

type TorchBuffer struct {
	cTorchBuffer C.struct_vaccel_torch_buffer
}

func (b *TorchBuffer) Init(data string) int {
	return int(C.vaccel_torch_buffer_init(&b.cTorchBuffer, C.CString(data), C.size_t(len(data)))) //nolint:gocritic
}

func (b *TorchBuffer) Release() int {
	return int(C.vaccel_torch_buffer_release(&b.cTorchBuffer)) //nolint:gocritic
}

func (b *TorchBuffer) TakeData() string {
	if b == nil {
		return ""
	}

	str := C.GoString(b.cTorchBuffer.data)

	b.cTorchBuffer.data = nil
	b.cTorchBuffer.size = 0

	return str
}

type TorchTensor struct {
	cTorchTensor *C.struct_vaccel_torch_tensor
}

func (t *TorchTensor) Init(dims []int64, dtype TorchDataType) int {
	if t == nil || len(dims) == 0 {
		return EINVAL
	}

	cDims := (*C.int64_t)(C.malloc(C.size_t(len(dims)) * C.size_t(unsafe.Sizeof(C.int64_t(0)))))
	if cDims == nil {
		return ENOMEM
	}
	defer C.free(unsafe.Pointer(cDims))

	base := uintptr(unsafe.Pointer(cDims))
	for i := 0; i < len(dims); i++ {
		ptr := (*C.int64_t)(unsafe.Pointer(base + uintptr(i)*unsafe.Sizeof(C.int64_t(0))))
		*ptr = C.int64_t(dims[i])
	}

	ret := C.vaccel_torch_tensor_new(
		&t.cTorchTensor,
		C.int64_t(len(dims)),
		cDims,
		C.enum_vaccel_torch_data_type(dtype), //nolint:gocritic
	)

	return int(ret)
}

func (t *TorchTensor) Release() int {
	return int(C.vaccel_torch_tensor_delete(t.cTorchTensor)) //nolint:gocritic
}

func (t *TorchTensor) Allocate(dims []int64, dtype TorchDataType, totalSize uint) int {
	ret := t.Init(dims, dtype)
	if ret != OK {
		return ret
	}

	if totalSize == 0 {
		return OK
	}

	t.cTorchTensor.data = C.malloc(C.size_t(totalSize))
	if t.cTorchTensor.data == nil {
		C.vaccel_torch_tensor_delete(t.cTorchTensor) //nolint:gocritic
		return ENOMEM
	}

	t.cTorchTensor.size = C.size_t(totalSize)
	t.cTorchTensor.owned = true

	return OK
}

func (t *TorchTensor) SetData(data uintptr, size uint, own bool) int {
	if t == nil {
		return EINVAL
	}

	if t.cTorchTensor.data != nil && t.cTorchTensor.owned {
		fmt.Println("Previous tensor data will not be freed by release!")
	}

	t.cTorchTensor.data = unsafe.Pointer(data)
	t.cTorchTensor.size = C.size_t(size)
	t.cTorchTensor.owned = C.bool(own)

	return OK
}

func (t *TorchTensor) TakeData() (uintptr, uint) {
	if t == nil {
		return 0, 0
	}

	outData := uintptr(t.cTorchTensor.data)
	outSize := uint(t.cTorchTensor.size)

	t.cTorchTensor.data = nil
	t.cTorchTensor.size = 0
	t.cTorchTensor.owned = false

	return outData, outSize
}

func (t *TorchTensor) Data() uintptr {
	if t == nil {
		return 0
	}
	return uintptr(t.cTorchTensor.data)
}

func (t *TorchTensor) NrDims() int {
	if t == nil || t.cTorchTensor.dims == nil || t.cTorchTensor.nr_dims <= 0 {
		return -1
	}
	return int(t.cTorchTensor.nr_dims)
}

func (t *TorchTensor) Type() TorchDataType {
	if t == nil || t.cTorchTensor.data_type < 1 {
		return -1
	}
	return TorchDataType(t.cTorchTensor.data_type)
}

func (t *TorchTensor) Dims() []int64 {
	if t == nil || t.cTorchTensor.dims == nil || t.cTorchTensor.nr_dims <= 0 {
		return nil
	}

	n := int(t.cTorchTensor.nr_dims)
	ptr := unsafe.Pointer(t.cTorchTensor.dims)
	cSlice := unsafe.Slice((*int64)(ptr), n)
	dims := make([]int64, n)
	copy(dims, cSlice)

	return dims
}

func (t *TorchTensor) DataPtr() uintptr {
	if t == nil || t.cTorchTensor.data == nil {
		return 0
	}
	return uintptr(t.cTorchTensor.data)
}

func TorchModelLoad(sess *Session, model *Resource) int {
	if sess == nil || model == nil {
		return EINVAL
	}

	return int(C.vaccel_torch_model_load(&sess.cSess, model.cRes)) //nolint:gocritic
}

func TorchModelRun(
	sess *Session,
	model *Resource,
	buffer *TorchBuffer,
	inTensors []TorchTensor,
	outTensors *[]TorchTensor,
) int {
	if sess == nil || model == nil || outTensors == nil {
		return EINVAL
	}

	nrInputs := len(inTensors)
	nrOutputs := len(*outTensors)
	if nrOutputs == 0 {
		return EINVAL
	}

	inBufSize := C.size_t(nrInputs) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cInPtr := C.malloc(inBufSize)
	defer C.free(cInPtr)

	inTensorSlice := unsafe.Slice((**C.struct_vaccel_torch_tensor)(cInPtr), nrInputs)
	for i := 0; i < nrInputs; i++ {
		inTensorSlice[i] = inTensors[i].cTorchTensor
	}

	outBufSize := C.size_t(nrOutputs) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cOutPtr := C.malloc(outBufSize)
	defer C.free(cOutPtr)

	var bufPtr *C.struct_vaccel_torch_buffer
	if buffer != nil {
		bufPtr = &buffer.cTorchBuffer
	} else {
		bufPtr = nil
	}

	ret := int(C.vaccel_torch_model_run(
		&sess.cSess,
		model.cRes,
		bufPtr,
		(**C.struct_vaccel_torch_tensor)(cInPtr),
		C.int(nrInputs),
		(**C.struct_vaccel_torch_tensor)(cOutPtr),
		C.int(nrOutputs), //nolint:gocritic
	))
	if ret != OK {
		return ret
	}

	outTensorSlice := unsafe.Slice((**C.struct_vaccel_torch_tensor)(cOutPtr), nrInputs)
	for i := 0; i < nrOutputs; i++ {
		if outTensorSlice[i] != nil {
			nrDims := int(outTensorSlice[i].nr_dims)
			dtype := TorchDataType(outTensorSlice[i].data_type)

			dims := make([]int64, nrDims)
			src := unsafe.Slice((*int64)(unsafe.Pointer(outTensorSlice[i].dims)), nrDims)
			copy(dims, src)

			ret = (*outTensors)[i].Init(dims, dtype)
			if ret != OK {
				fmt.Println("Could not initialize output tensor from C dims and type")
				return ret
			}

			var data unsafe.Pointer
			var size C.size_t
			ret = int(C.vaccel_torch_tensor_take_data(outTensorSlice[i], &data, &size)) //nolint:gocritic
			if ret != OK {
				fmt.Println("Could not take data from output C tensor")
				return ret
			}

			ret = int(C.vaccel_torch_tensor_delete(outTensorSlice[i]))
			if ret != OK {
				fmt.Println("Could not delete output C tensor")
				return ret
			}

			goData := uintptr(data)
			goSize := uint(size)

			ret = (*outTensors)[i].SetData(goData, goSize, true)
			if ret != OK {
				fmt.Println("Could not set data for Golang output tensor")
				return ret
			}
		}
	}

	return OK
}
