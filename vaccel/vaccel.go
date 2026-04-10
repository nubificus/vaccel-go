// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl

// TODO: Remove this once deprecated functions are updated
#cgo CFLAGS: -Wno-deprecated -Wno-deprecated-declarations

#include <vaccel.h>
*/
import "C"
