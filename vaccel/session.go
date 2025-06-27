// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#cgo pkg-config: vaccel
#cgo LDFLAGS: -lvaccel -ldl
#include <vaccel.h>

*/
import "C"

type Session struct {
	cSess C.struct_vaccel_session
}

func (s *Session) Init(flags uint32) int {
	return int(C.vaccel_session_init(&s.cSess, C.uint32_t(flags))) //nolint:gocritic
}

func (s *Session) Release() int {
	return int(C.vaccel_session_release(&s.cSess)) //nolint:gocritic
}

func (s *Session) Register(r *Resource) int {
	return int(C.vaccel_resource_register(&r.cRes, &s.cSess)) //nolint:gocritic
}

func (s *Session) Unregister(r *Resource) int {
	return int(C.vaccel_resource_unregister(&r.cRes, &s.cSess)) //nolint:gocritic
}

func (s *Session) GetId() int64 {
	return int64(s.cSess.id)
}
