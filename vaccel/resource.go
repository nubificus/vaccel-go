// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"

type ResourceType int

const (
	ResourceLib ResourceType = iota
	ResourceData
	ResourceModel
)

type Resource struct {
	cRes C.struct_vaccel_resource
}

func (t ResourceType) ToCEnum() C.vaccel_resource_type_t {
	return C.vaccel_resource_type_t(t)
}

func (r *Resource) Init(path string, resType ResourceType) int {
	return int(C.vaccel_resource_init(&r.cRes, C.CString(path), resType.ToCEnum())) //nolint:gocritic
}

func (r *Resource) Release() int {
	return int(C.vaccel_resource_release(&r.cRes)) //nolint:gocritic
}
