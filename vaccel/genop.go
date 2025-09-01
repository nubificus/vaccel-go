// SPDX-License-Identifier: Apache-2.0

package vaccel

/*

#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"

func Genop(sess *Session, read *ArgList, write *ArgList) int {
	cRead := read.cList.list
	cWrite := write.cList.list

	cNrRead := C.int(read.cList.size)
	cNrWrite := C.int(write.cList.size)

	cRet := C.vaccel_genop(sess.cSess, cRead, cNrRead, cWrite, cNrWrite) //nolint:gocritic
	return int(cRet)

}
