// SPDX-License-Identifier: Apache-2.0

package vaccel

// #include <vaccel/ops/noop.h>
import "C"

func NoOp(sess *Session) int {
	return int(C.vaccel_noop(sess.cSess))
}
