// SPDX-License-Identifier: Apache-2.0

package vaccel

// #include <vaccel/session.h>
import "C"

type Session struct {
	cSess *C.struct_vaccel_session
}

func (s *Session) Init(flags uint32) int {
	return int(C.vaccel_session_new(&s.cSess, C.uint32_t(flags)))
}

func (s *Session) Release() int {
	return int(C.vaccel_session_delete(s.cSess))
}

func (s *Session) Register(r *Resource) int {
	return int(C.vaccel_resource_register(r.cRes, s.cSess))
}

func (s *Session) Unregister(r *Resource) int {
	return int(C.vaccel_resource_unregister(r.cRes, s.cSess))
}

func (s *Session) GetID() int64 {
	return int64(s.cSess.id)
}

func (s *Session) Update(flags uint32) int {
	return int(C.vaccel_session_update(s.cSess, C.uint32_t(flags)))
}

func (s *Session) GetFlags() int32 {
	return int32(s.cSess.hint)
}
