// SPDX-License-Identifier: Apache-2.0

package vaccel

// #include <stdlib.h>
// #include <string.h>
// #include <vaccel/arg.h>
import "C"
import "unsafe"

type Arg struct {
	cArg *C.struct_vaccel_arg
}

type ArgList struct {
	cList *C.struct_vaccel_arg_list
}

/* Type of function to serialize a structure */
/* Returns pointer to serialized data and the size in bytes */
type Serializer func(buf unsafe.Pointer) (unsafe.Pointer, uint32)

/* Type of function to deserialize a structure */
/* Returns pointer to the constructed structure */
type Deserializer func(buf unsafe.Pointer) unsafe.Pointer

// Deprecated: The C vaccel_arg_list API is deprecated
func ArgsInit(size uint32) *ArgList {
	list := new(ArgList)
	list.cList = C.vaccel_args_init(C.uint(size))

	return list
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) AddSerialArg(buf unsafe.Pointer, size int) int {
	return int(C.vaccel_add_serial_arg(arglist.cList, buf, C.uint(size)))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) AddStringArg(arg string) int {
	cStr := C.CString(arg)
	ret := int(C.vaccel_add_serial_arg(arglist.cList, unsafe.Pointer(cStr), C.uint(C.strlen(cStr))))
	if ret == OK {
		length := int(arglist.cList.size)
		idx := int(arglist.cList.curr_idx - 1)
		slice := (*[1 << 16]C.int)(unsafe.Pointer(arglist.cList.idcs_allocated_space))[:length:length]
		slice[idx] = 1
	} else {
		C.free(unsafe.Pointer(cStr))
	}
	return ret
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) AddInt32Arg(arg int32) int {
	cInt := (*C.int32_t)(C.malloc(C.sizeof_int32_t))
	*cInt = C.int32_t(arg)
	ret := int(C.vaccel_add_serial_arg(arglist.cList, unsafe.Pointer(cInt), C.sizeof_int32_t))
	if ret == OK {
		length := int(arglist.cList.size)
		idx := int(arglist.cList.curr_idx - 1)
		slice := (*[1 << 16]C.int)(unsafe.Pointer(arglist.cList.idcs_allocated_space))[:length:length]
		slice[idx] = 1
	} else {
		C.free(unsafe.Pointer(cInt))
	}
	return ret
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) AddNonSerialArg(nonSerialBuf unsafe.Pointer,
	argtype uint32, serialize Serializer) int { //nolint:revive // argtype will be used in a next iteration
	serialBuf, bytes := serialize(nonSerialBuf)
	return arglist.AddSerialArg(serialBuf, int(bytes))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) ExpectSerialArg(buf unsafe.Pointer, size int) int {
	return int(C.vaccel_expect_serial_arg(arglist.cList, buf, C.uint(size)))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) ExpectNonSerialArg(expectedSize int) int {
	return int(C.vaccel_expect_nonserial_arg(arglist.cList, C.uint(expectedSize)))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) GetArgs() *Arg {
	args := new(Arg)
	args.cArg = arglist.cList.list
	return args
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (args *Arg) ExtractSerialArg(idx int) unsafe.Pointer {
	return C.vaccel_extract_serial_arg(args.cArg, C.int(idx))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) ExtractSerialArg(idx int) unsafe.Pointer {
	return C.vaccel_extract_serial_arg(arglist.cList.list, C.int(idx))
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) ExtractNonSerialArg(idx int, deserialize Deserializer) unsafe.Pointer {
	nonSerialBuf := arglist.ExtractSerialArg(idx)
	return deserialize(nonSerialBuf)
}

// Deprecated: The C vaccel_arg_list API is deprecated
func (arglist *ArgList) Delete() int {
	return int(C.vaccel_delete_args(arglist.cList))
}
