// SPDX-License-Identifier: Apache-2.0

package vaccel

/*

#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"

func ExecWithResource(sess *Session, res *Resource, funcname string,
	read *ArgList, write *ArgList) int {
	cfunc := C.CString(funcname)
	cread := read.cList.list
	cwrite := write.cList.list

	cNrRead := C.size_t(read.cList.size)
	cNrWrite := C.size_t(write.cList.size)

	cRet := C.vaccel_exec_with_resource(sess.cSess, res.cRes, cfunc, cread, cNrRead, cwrite, cNrWrite) //nolint:gocritic

	return int(cRet)

}
