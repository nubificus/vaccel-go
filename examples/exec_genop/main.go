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

const INPUT int32 = 10

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s <lib_file> [iterations]\n", os.Args[0])
		return
	}

	libpath := os.Args[1]
	funcname := "mytestfunc"
	iters := 1
	var e error
	if len(os.Args) > 2 {
		iters, e = strconv.Atoi(os.Args[2])
		if e != nil {
			fmt.Println("error converting input iterations, assuming 1..")
		}
	}

	var output int32
	var session vaccel.Session
	var read, write *vaccel.ArgList

	err := session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(int(err))
	}

	read = vaccel.ArgsInit(4)
	if read == nil {
		fmt.Println("Error Creating the read arg-list")
		goto ReleaseSession
	}

	if read.AddInt32Arg(int32(vaccel.OpExec)) != vaccel.OK {
		fmt.Println("Error Adding Serialized arg (type)")
		goto DeleteReadArgs
	}

	if read.AddStringArg(libpath) != vaccel.OK {
		fmt.Println("Error Adding Serialized arg (libPath)")
		goto DeleteReadArgs
	}

	if read.AddStringArg(funcname) != vaccel.OK {
		fmt.Println("Error Adding Serialized arg (func)")
		goto DeleteReadArgs
	}

	if read.AddInt32Arg(INPUT) != vaccel.OK {
		fmt.Println("Error Adding Serialized arg (input)")
		goto DeleteReadArgs
	}

	write = vaccel.ArgsInit(1)
	if write == nil {
		fmt.Println("Error Creating the write arg-list")
		goto DeleteReadArgs
	}

	if write.ExpectSerialArg(unsafe.Pointer(&output), 4) != vaccel.OK {
		fmt.Println("Error defining expected arg")
		goto DeleteWriteArgs
	}

	for i := 0; i < iters; i++ {
		err = vaccel.Genop(&session, read, write)
		if err != vaccel.OK {
			fmt.Println("An error occurred while running the operation")
			goto DeleteWriteArgs
		}
	}

	fmt.Println("Output: ", output)

DeleteWriteArgs:
	if write.Delete() != vaccel.OK {
		fmt.Println("An error occurred in deletion of the write arguments")
	}
DeleteReadArgs:
	if read.Delete() != vaccel.OK {
		fmt.Println("An error occurred in deletion of the read arguments")
	}
ReleaseSession:
	if session.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the session")
	}
}
