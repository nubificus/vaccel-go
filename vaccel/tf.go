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

type TFDataType int

const (
	TfFloat      TFDataType = 1
	TfDouble     TFDataType = 2
	TfInt32      TFDataType = 3
	TfUint8      TFDataType = 4
	TfInt16      TFDataType = 5
	TfInt8       TFDataType = 6
	TfString     TFDataType = 7
	TfComplex64  TFDataType = 8
	TfComplex    TFDataType = 8
	TfInt64      TFDataType = 9
	TfBool       TFDataType = 10
	TfQint8      TFDataType = 11
	TfQuint8     TFDataType = 12
	TfQint32     TFDataType = 13
	TfBfloat16   TFDataType = 14
	TfQint16     TFDataType = 15
	TfQuint16    TFDataType = 16
	TfUint16     TFDataType = 17
	TfComplex128 TFDataType = 18
	TfHalf       TFDataType = 19
	TfResource   TFDataType = 20
	TfVariant    TFDataType = 21
	TfUint32     TFDataType = 22
	TfUint64     TFDataType = 23
)

type TFBuffer struct {
	cTFBuf C.struct_vaccel_tf_buffer
}

func (b *TFBuffer) Init(data uintptr, size uint) int {
	cData := unsafe.Pointer(data)
	cSize := C.size_t(size)
	return int(C.vaccel_tf_buffer_init(&b.cTFBuf, cData, cSize)) //nolint:gocritic
}

func (b *TFBuffer) Release() int {
	return int(C.vaccel_tf_buffer_release(&b.cTFBuf)) //nolint:gocritic
}

func (b *TFBuffer) TakeData() (uintptr, uint) {
	if b == nil {
		return 0, 0
	}

	outData := uintptr(b.cTFBuf.data)
	outSize := uint(b.cTFBuf.size)

	b.cTFBuf.data = nil
	b.cTFBuf.size = 0

	return outData, outSize
}

type TFNode struct {
	cTFNode C.struct_vaccel_tf_node
}

func (n *TFNode) Init(name string, id int) int {
	cInt := C.int(id)
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.vaccel_tf_node_init(&n.cTFNode, cStr, cInt)) //nolint:gocritic
}

func (n *TFNode) Release() int {
	return int(C.vaccel_tf_node_release(&n.cTFNode)) //nolint:gocritic
}

type TFTensor struct {
	cTFTensor C.struct_vaccel_tf_tensor
}

func (t *TFTensor) Init(dims []int64, dtype TFDataType) int {
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

	ret := C.vaccel_tf_tensor_init(
		&t.cTFTensor,
		C.int(len(dims)),
		cDims,
		C.enum_vaccel_tf_data_type(dtype), //nolint:gocritic
	)

	return int(ret)
}

func (t *TFTensor) Release() int {
	return int(C.vaccel_tf_tensor_release(&t.cTFTensor)) //nolint:gocritic
}

func (t *TFTensor) Allocate(dims []int64, dtype TFDataType, totalSize uint) int {
	ret := t.Init(dims, dtype)
	if ret != OK {
		return ret
	}

	if totalSize == 0 {
		return OK
	}

	t.cTFTensor.data = C.malloc(C.size_t(totalSize))
	if t.cTFTensor.data == nil {
		C.vaccel_tf_tensor_release(&t.cTFTensor) //nolint:gocritic
		return ENOMEM
	}

	t.cTFTensor.size = C.size_t(totalSize)
	t.cTFTensor.owned = true

	return OK
}

func (t *TFTensor) SetData(data uintptr, size uint, own bool) int {
	if t == nil {
		return EINVAL
	}

	if t.cTFTensor.data != nil && t.cTFTensor.owned {
		fmt.Println("Previous tensor data will not be freed by release!")
	}

	t.cTFTensor.data = unsafe.Pointer(data)
	t.cTFTensor.size = C.size_t(size)
	t.cTFTensor.owned = C.bool(own)

	return OK
}

func (t *TFTensor) TakeData() (uintptr, uint) {
	if t == nil {
		return 0, 0
	}

	outData := uintptr(t.cTFTensor.data)
	outSize := uint(t.cTFTensor.size)

	t.cTFTensor.data = nil
	t.cTFTensor.size = 0
	t.cTFTensor.owned = false

	return outData, outSize
}

func (t *TFTensor) Data() uintptr {
	if t == nil {
		return 0
	}
	return uintptr(t.cTFTensor.data)
}

func (t *TFTensor) NrDims() int {
	if t == nil || t.cTFTensor.dims == nil || t.cTFTensor.nr_dims <= 0 {
		return -1
	}
	return int(t.cTFTensor.nr_dims)
}

func (t *TFTensor) Type() TFDataType {
	if t == nil || t.cTFTensor.data_type < 1 {
		return -1
	}
	return TFDataType(t.cTFTensor.data_type)
}

func (t *TFTensor) Dims() []int64 {
	if t == nil || t.cTFTensor.dims == nil || t.cTFTensor.nr_dims <= 0 {
		return nil
	}

	n := int(t.cTFTensor.nr_dims)
	ptr := unsafe.Pointer(t.cTFTensor.dims)
	cSlice := unsafe.Slice((*int64)(ptr), n)
	dims := make([]int64, n)
	copy(dims, cSlice)

	return dims
}

func (t *TFTensor) DataPtr() uintptr {
	if t == nil || t.cTFTensor.data == nil {
		return 0
	}
	return uintptr(t.cTFTensor.data)
}

func (t *TFTensor) PrintFloat32Data() {
	if t == nil || t.cTFTensor.data == nil {
		fmt.Println("nil tensor")
		return
	}

	if t.cTFTensor.data_type != C.VACCEL_TF_FLOAT {
		fmt.Println("Unsupported data type (only float32 supported)")
		return
	}

	dims := t.Dims()
	if dims == nil {
		fmt.Println("Invalid dims")
		return
	}

	numel := int64(1)
	for _, d := range dims {
		numel *= d
	}

	ptr := (*float32)(unsafe.Pointer(t.DataPtr()))
	slice := unsafe.Slice(ptr, numel)

	fmt.Printf("Tensor shape: %v\n", dims)
	fmt.Println("Values:")
	printRecursiveFloat32(slice, dims, 0)
}

func printRecursiveFloat32(data []float32, dims []int64, level int) {
	if len(dims) == 0 {
		return
	}
	if len(dims) == 1 {
		fmt.Printf("%s[", indent(level))
		for i := int64(0); i < dims[0]; i++ {
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
	for i := int64(0); i < size; i++ {
		printRecursiveFloat32(data[i*int64(sub):(i+1)*int64(sub)], dims[1:], level+1)
	}
}

func indent(level int) string {
	return "  " + string(make([]rune, level*2))
}

type TFStatus struct {
	cTFStatus C.struct_vaccel_tf_status
}

func (s *TFStatus) Init(errorCode uint8, message string) int {
	cErr := C.uint8_t(errorCode)
	cMsg := C.CString(message)
	defer C.free(unsafe.Pointer(cMsg))
	return int(C.vaccel_tf_status_init(&s.cTFStatus, cErr, cMsg)) //nolint:gocritic
}

func (s *TFStatus) Release() int {
	if s == nil {
		return EINVAL
	}
	return int(C.vaccel_tf_status_release(&s.cTFStatus)) //nolint:gocritic
}

func TFModelLoad(sess *Session, model *Resource, status *TFStatus) int {
	if sess == nil || model == nil {
		return EINVAL
	}
	return int(C.vaccel_tf_model_load(&sess.cSess, &model.cRes, &status.cTFStatus)) //nolint:gocritic
}

func TFModelRun(
	sess *Session,
	model *Resource,
	runOptions *TFBuffer,
	inNodes *TFNode,
	inTensors []TFTensor,
	outNodes *TFNode,
	outTensors *[]TFTensor,
	status *TFStatus,
) int {
	if sess == nil || model == nil || inNodes == nil || outNodes == nil || status == nil || outTensors == nil {
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

	inTensorSlice := unsafe.Slice((**C.struct_vaccel_tf_tensor)(cInPtr), nrInputs)
	for i := 0; i < nrInputs; i++ {
		inTensorSlice[i] = &inTensors[i].cTFTensor
	}

	outBufSize := C.size_t(nrOutputs) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cOutPtr := C.malloc(outBufSize)
	defer C.free(cOutPtr)

	ret := int(C.vaccel_tf_model_run(
		&sess.cSess,
		&model.cRes,
		func() *C.struct_vaccel_tf_buffer {
			if runOptions != nil {
				return &runOptions.cTFBuf
			}
			return nil
		}(),
		&inNodes.cTFNode,
		(**C.struct_vaccel_tf_tensor)(cInPtr),
		C.int(nrInputs),
		&outNodes.cTFNode,
		(**C.struct_vaccel_tf_tensor)(cOutPtr),
		C.int(nrOutputs),
		&status.cTFStatus, //nolint:gocritic
	))
	if ret != OK {
		return ret
	}

	outTensorSlice := unsafe.Slice((**C.struct_vaccel_tf_tensor)(cOutPtr), nrInputs)
	for i := 0; i < nrOutputs; i++ {
		if outTensorSlice[i] != nil {
			nrDims := int(outTensorSlice[i].nr_dims)
			dtype := TFDataType(outTensorSlice[i].data_type)

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
			ret = int(C.vaccel_tf_tensor_take_data(outTensorSlice[i], &data, &size)) //nolint:gocritic
			if ret != OK {
				fmt.Println("Could not take data from output C tensor")
				return ret
			}

			ret = int(C.vaccel_tf_tensor_delete(outTensorSlice[i]))
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

func TFModelUnload(sess *Session, model *Resource, status *TFStatus) int {
	err := int(C.vaccel_tf_model_unload(&sess.cSess, &model.cRes, &status.cTFStatus)) //nolint:gocritic
	return err
}
