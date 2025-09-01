// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"

func NoOp(sess *Session) int {
	return int(C.vaccel_noop(sess.cSess)) //nolint:gocritic
}
