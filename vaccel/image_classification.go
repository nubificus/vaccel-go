// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func ImageClassificationFromFile(sess *Session, imagePath string) (string, int) {

	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	cImageBytes := (*C.uchar)(&imageBytes[0])
	cImgBuf := unsafe.Pointer(cImageBytes)
	cImgLen := C.size_t(len(imageBytes))

	cText := (*C.uchar)(C.malloc(C.size_t(256)))
	cOutImageName := (*C.uchar)(C.malloc(C.size_t(256)))

	/* Free the memory when done */
	defer C.free(unsafe.Pointer(cText))
	defer C.free(unsafe.Pointer(cOutImageName))

	cRet := C.vaccel_image_classification(
		sess.cSess, cImgBuf, cText, cOutImageName,
		cImgLen, C.size_t(256), C.size_t(256)) //nolint:gocritic

	var golangOut string

	if int(cRet) == 0 {
		ptr := unsafe.Pointer(cText)
		typeCast := (*C.char)(ptr)
		golangOut = C.GoString(typeCast)
	} else {
		golangOut =
			"A problem occurred while running the Operation"
	}

	return golangOut, int(cRet)
}

func ImageClassification(sess *Session, image []byte) (string, int) {

	cImageBytes := (*C.uchar)(&image[0])
	cImgBuf := unsafe.Pointer(cImageBytes)
	cImgLen := C.size_t(len(image))

	cText := (*C.uchar)(C.malloc(C.size_t(256)))
	cOutImageName := (*C.uchar)(C.malloc(C.size_t(256)))

	/* Free the memory when done */
	defer C.free(unsafe.Pointer(cText))
	defer C.free(unsafe.Pointer(cOutImageName))

	cRet := C.vaccel_image_classification(
		sess.cSess, cImgBuf, cText, cOutImageName,
		cImgLen, C.size_t(256), C.size_t(256)) //nolint:gocritic

	var golangOut string

	if int(cRet) == 0 {
		ptr := unsafe.Pointer(cText)
		typeCast := (*C.char)(ptr)
		golangOut = C.GoString(typeCast)
	} else {
		golangOut =
			"A problem occurred while running the Operation"
	}

	return golangOut, int(cRet)
}
