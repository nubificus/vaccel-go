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

func BufferInit(buffer *TFBuffer, data uintptr, size uint) int {
	cData := unsafe.Pointer(data)
	cSize := C.size_t(size)
	return int(C.vaccel_tf_buffer_init(&buffer.cTFBuf, cData, cSize)) //nolint:gocritic
}

func TFBufferRelease(buffer *TFBuffer) int {
	return int(C.vaccel_tf_buffer_release(&buffer.cTFBuf)) //nolint:gocritic
}

func TFBufferTakeData(buffer *TFBuffer, outData *uintptr, outSize *uint) int {
	if buffer == nil || outData == nil || outSize == nil {
		return EINVAL
	}

	*outData = uintptr(buffer.cTFBuf.data)
	*outSize = uint(buffer.cTFBuf.size)

	buffer.cTFBuf.data = nil
	buffer.cTFBuf.size = 0

	return OK
}

type TFNode struct {
	cTFNode C.struct_vaccel_tf_node
}

func TFNodeInit(node *TFNode, name string, id int) int {
	cInt := C.int(id)
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.vaccel_tf_node_init(&node.cTFNode, cStr, cInt)) //nolint:gocritic
}

func TFNodeRelease(node *TFNode) int {
	return int(C.vaccel_tf_node_release(&node.cTFNode)) //nolint:gocritic
}

type TFTensor struct {
	cTFTensor C.struct_vaccel_tf_tensor
}

func TFTensorInit(tensor *TFTensor, dims []int64, dtype TFDataType) int {
	if tensor == nil || len(dims) == 0 {
		return EINVAL
	}

	cDims := (*C.int64_t)(C.malloc(C.size_t(len(dims)) * C.size_t(unsafe.Sizeof(C.int64_t(0)))))
	if cDims == nil {
		return ENOMEM
	}
	defer C.free(unsafe.Pointer(cDims))

	dimsSlice := (*[16]C.int64_t)(unsafe.Pointer(cDims))[:len(dims):len(dims)]
	for i, d := range dims {
		dimsSlice[i] = C.int64_t(d)
	}

	ret := C.vaccel_tf_tensor_init(
		&tensor.cTFTensor,
		C.int(len(dims)),
		cDims,
		C.enum_vaccel_tf_data_type(dtype), //nolint:gocritic
	)

	return int(ret)
}

func TFTensorRelease(tensor *TFTensor) int {
	return int(C.vaccel_tf_tensor_release(&tensor.cTFTensor)) //nolint:gocritic
}

func TFTensorAllocate(tensor *TFTensor, dims []int64, dtype TFDataType, totalSize uint) int {
	ret := TFTensorInit(tensor, dims, dtype)
	if ret != OK {
		return ret
	}

	if totalSize == 0 {
		return OK
	}

	tensor.cTFTensor.data = C.malloc(C.size_t(totalSize))
	if tensor.cTFTensor.data == nil {
		C.vaccel_tf_tensor_release(&tensor.cTFTensor) //nolint:gocritic
		return ENOMEM
	}

	tensor.cTFTensor.size = C.size_t(totalSize)
	tensor.cTFTensor.owned = true

	return OK
}

func TFTensorSetData(tensor *TFTensor, data uintptr, size uint, own bool) int {
	if tensor == nil {
		return EINVAL
	}

	if tensor.cTFTensor.data != nil && tensor.cTFTensor.owned {
		fmt.Println("Previous tensor data will not be freed by release!")
	}

	tensor.cTFTensor.data = unsafe.Pointer(data)
	tensor.cTFTensor.size = C.size_t(size)
	tensor.cTFTensor.owned = C.bool(own)

	return OK
}

func TFTensorTakeData(tensor *TFTensor, outData *uintptr, outSize *uint) int {
	if tensor == nil || outData == nil || outSize == nil {
		return EINVAL
	}

	*outData = uintptr(tensor.cTFTensor.data)
	*outSize = uint(tensor.cTFTensor.size)

	tensor.cTFTensor.data = nil
	tensor.cTFTensor.size = 0
	tensor.cTFTensor.owned = false

	return OK
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

func (t *TFTensor) Type() int {
	if t == nil || t.cTFTensor.data_type < 1 {
		return -1
	}
	return int(t.cTFTensor.data_type)
}

func (t *TFTensor) Dims() []int64 {
	if t == nil || t.cTFTensor.dims == nil || t.cTFTensor.nr_dims <= 0 {
		return nil
	}

	n := int(t.cTFTensor.nr_dims)
	cDims := (*[16]C.int64_t)(unsafe.Pointer(t.cTFTensor.dims))[:n:n]

	dims := make([]int64, n)
	for i := 0; i < n; i++ {
		dims[i] = int64(cDims[i])
	}
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

func TFStatusInit(status *TFStatus, errorCode uint8, message string) int {
	if status == nil {
		return EINVAL
	}

	cMsg := C.CString(message)
	defer C.free(unsafe.Pointer(cMsg))

	ret := C.vaccel_tf_status_init(&status.cTFStatus, C.uint8_t(errorCode), cMsg) //nolint:gocritic
	return int(ret)
}

func TFStatusRelease(status *TFStatus) int {
	if status == nil {
		return EINVAL
	}
	return int(C.vaccel_tf_status_release(&status.cTFStatus)) //nolint:gocritic
}

func TFSessionLoad(sess *Session, model *Resource, status *TFStatus) int {
	return int(C.vaccel_tf_session_load(&sess.cSess, &model.cRes, &status.cTFStatus)) //nolint:gocritic
}

func TFSessionRun(
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

	ret := int(C.vaccel_tf_session_run(
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

			ret = TFTensorInit(&(*outTensors)[i], dims, dtype)
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

			ret = TFTensorSetData(&(*outTensors)[i], goData, goSize, true)
			if ret != OK {
				fmt.Println("Could not set data for Golang output tensor")
				return ret
			}
		}
	}

	return OK
}

func TFSessionDelete(sess *Session, model *Resource, status *TFStatus) int {
	return int(C.vaccel_tf_session_delete(&sess.cSess, &model.cRes, &status.cTFStatus)) //nolint:gocritic
}
