// SPDX-License-Identifier: Apache-2.0

package main

import (
	"C"
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {

	/* Read User Args */
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./exec <filename> <input>")
		return
	}

	path := os.Args[1]
	input := os.Args[2]

	inputInt, e := strconv.Atoi(input)
	if e != nil {
		fmt.Println("error converting input")
	}

	var session vaccel.Session
	err := session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(int(err))
	}

	var res vaccel.Resource
	err = res.Init(path, vaccel.ResourceLib)
	if err != vaccel.OK {
		fmt.Println("error creating shared object resource")
		os.Exit(int(err))
	}

	err = session.Register(&res)
	if err != vaccel.OK {
		fmt.Println("error registering resource with session")
		os.Exit(int(err))
	}

	read := vaccel.ArgsInit(1)
	write := vaccel.ArgsInit(1)

	if read == nil || write == nil {
		fmt.Println("Error Creating the arg-lists")
		os.Exit(0)
	}

	buf := unsafe.Pointer(&inputInt)
	size := unsafe.Sizeof(inputInt)

	if read.AddSerialArg(buf, size) != vaccel.OK {
		fmt.Println("Error Adding Serialized arg")
		os.Exit(0)
	}

	var output int
	buf = unsafe.Pointer(&output)
	size = unsafe.Sizeof(output)

	if write.ExpectSerialArg(buf, size) != vaccel.OK {
		fmt.Println("Error defining expected arg")
		os.Exit(0)
	}

	err = vaccel.ExecWithResource(&session, &res, "mytestfunc", read, write)
	if err != vaccel.OK {
		fmt.Println("An error occurred while running the operation")
		os.Exit(err)
	}

	fmt.Println("Output(1): ", C.uint(output))

	/* Or */
	outbuf := write.ExtractSerialArg(0)

	cast := (*int)(outbuf)
	val := C.uint(*cast)
	fmt.Println("Output(2): ", val)

	/* Or */
	outbuf = write.GetArgs().ExtractSerialArg(0)

	cast = (*int)(outbuf)
	val = C.uint(*cast)
	fmt.Println("Output(3): ", val)

	if write.Delete() != vaccel.OK || read.Delete() != vaccel.OK {
		fmt.Println("An error occurred in deletion of the arg-lists")
		os.Exit(0)
	}

	if session.Unregister(&res) != vaccel.OK {
		fmt.Println("An error occurred while unregistering the resource")
	}

	if res.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the resource")
	}

	if session.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the session")
	}
}
