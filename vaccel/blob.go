// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"
import "unsafe"

type BlobType int32

const (
	BlobFile   BlobType = C.VACCEL_BLOB_FILE
	BlobBuffer BlobType = C.VACCEL_BLOB_BUFFER
	BlobMapped BlobType = C.VACCEL_BLOB_MAPPED
)

type Blob struct {
	cBlob *C.struct_vaccel_blob
}

func (b *Blob) Init(path string) int {
	return int(C.vaccel_blob_new(&b.cBlob, C.CString(path))) //nolint:gocritic
}

func (b *Blob) InitFromBuf(bytes []byte, own bool, filename string, dir string, randomize bool) int {
	var cdname *C.char

	cBlobBytes := (*C.uchar)(&bytes[0])
	cBlobLen := C.size_t(len(bytes))

	if dir == "" {
		cdname = nil
	} else {
		cdname = C.CString(dir)
		defer C.free(unsafe.Pointer(cdname))
	}

	return int(C.vaccel_blob_from_buf(&b.cBlob, cBlobBytes, cBlobLen, C.bool(own), C.CString(filename), cdname, C.bool(randomize))) //nolint:gocritic
}

func (b *Blob) Release() int {
	return int(C.vaccel_blob_delete(b.cBlob)) //nolint:gocritic
}
