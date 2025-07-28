// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>
#include <stdatomic.h>
#include <stdint.h>

static uint32_t get_refcount(atomic_uint *ref) {
	return atomic_load((_Atomic unsigned int *)ref);
}

*/
import "C"
import "unsafe"

type ResourceType int

const (
	ResourceLib ResourceType = iota
	ResourceData
	ResourceModel
)

type Resource struct {
	cRes *C.struct_vaccel_resource
}

func (t ResourceType) ToCEnum() C.vaccel_resource_type_t {
	return C.vaccel_resource_type_t(t)
}

func (r *Resource) Init(path string, resType ResourceType) int {
	return int(C.vaccel_resource_new(&r.cRes, C.CString(path), resType.ToCEnum())) //nolint:gocritic
}

func (r *Resource) InitMulti(paths []string, resType ResourceType) int {
	nrPaths := len(paths)
	cNrPaths := C.size_t(nrPaths)

	spaceSize := cNrPaths * C.size_t(unsafe.Sizeof(uintptr(0)))
	cSpace := C.malloc(spaceSize)
	defer C.free(cSpace)

	cPathsPtr := (**C.char)(cSpace)
	pathSlice := unsafe.Slice((**C.char)(cPathsPtr), nrPaths)

	for i := 0; i < nrPaths; i++ {
		str := C.CString(paths[i])
		defer C.free(unsafe.Pointer(str))

		pathSlice[i] = str
	}
	return int(C.vaccel_resource_multi_new(&r.cRes, cPathsPtr, cNrPaths, resType.ToCEnum())) //nolint:gocritic
}

func (r *Resource) InitFromBuf(bytes []byte, resType ResourceType, filename string, memOnly bool) int {
	cResBytes := (*C.uchar)(&bytes[0])
	cResBuf := unsafe.Pointer(cResBytes)
	cResLen := C.size_t(len(bytes))

	var cfname *C.char
	if filename == "" {
		cfname = nil
	} else {
		cfname = C.CString(filename)
		defer C.free(unsafe.Pointer(cfname))
	}

	return int(C.vaccel_resource_from_buf(&r.cRes, cResBuf, cResLen, resType.ToCEnum(), cfname, C.bool(memOnly))) //nolint:gocritic
}

func (r *Resource) InitFromBlobs(blobs []Blob, resType ResourceType) int {
	nrBlobs := len(blobs)
	cNrBlobs := C.size_t(nrBlobs)

	bufSize := cNrBlobs * C.size_t(unsafe.Sizeof(uintptr(0)))
	cSpace := C.malloc(bufSize)
	defer C.free(cSpace)

	cBlobsPtr := (**C.struct_vaccel_blob)(cSpace)
	blobSlice := unsafe.Slice((**C.struct_vaccel_blob)(cBlobsPtr), nrBlobs)

	for i := 0; i < nrBlobs; i++ {
		blobSlice[i] = blobs[i].cBlob
	}

	return int(C.vaccel_resource_from_blobs(&r.cRes, cBlobsPtr, cNrBlobs, resType.ToCEnum())) //nolint:gocritic
}

func (r *Resource) Release() int {
	return int(C.vaccel_resource_delete(r.cRes)) //nolint:gocritic
}

func (r *Resource) GetID() int64 {
	return int64(r.cRes.id)
}

func (r *Resource) GetRefcount() uint32 {
	return uint32(C.get_refcount(&r.cRes.refcount))
}
