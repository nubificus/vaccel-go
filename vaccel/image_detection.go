// SPDX-License-Identifier: Apache-2.0

package vaccel

// #include <vaccel/ops/image.h>
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func ImageDetectionFromFile(sess *Session, imagePath string) (string, int) {

	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	cImageBytes := (*C.uchar)(&imageBytes[0])
	cImgBuf := unsafe.Pointer(cImageBytes)
	cImgLen := C.size_t(len(imageBytes))

	cOutImageName := (*C.uchar)(C.malloc(C.size_t(1024)))

	/* Free the memory when done */
	defer C.free(unsafe.Pointer(cOutImageName))

	cRet := C.vaccel_image_detection(
		sess.cSess, cImgBuf, cOutImageName,
		cImgLen, C.size_t(1024)) //nolint:gocritic

	var golangOut string

	if int(cRet) == 0 {
		ptr := unsafe.Pointer(cOutImageName)
		typeCast := (*C.char)(ptr)
		golangOut = C.GoString(typeCast)
	} else {
		golangOut =
			"A problem occurred while running the Operation"
	}

	return golangOut, int(cRet)
}

func ImageDetection(sess *Session, image []byte) (string, int) {

	cImageBytes := (*C.uchar)(&image[0])
	cImgBuf := unsafe.Pointer(cImageBytes)
	cImgLen := C.size_t(len(image))

	cOutImageName := (*C.uchar)(C.malloc(C.size_t(1024)))

	/* Free the memory when done */
	defer C.free(unsafe.Pointer(cOutImageName))

	cRet := C.vaccel_image_detection(
		sess.cSess, cImgBuf, cOutImageName,
		cImgLen, C.size_t(1024)) //nolint:gocritic

	var golangOut string

	if int(cRet) == 0 {
		ptr := unsafe.Pointer(cOutImageName)
		typeCast := (*C.char)(ptr)
		golangOut = C.GoString(typeCast)
	} else {
		golangOut =
			"A problem occurred while running the Operation"
	}

	return golangOut, int(cRet)
}
