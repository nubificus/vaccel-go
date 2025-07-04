// SPDX-License-Identifier: Apache-2.0

package main

import (
	"C"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/nubificus/vaccel-go/vaccel"
)

type MyData struct {
	Size uint32
	Arr  []uint32
}

func NewMyData(size uint32) unsafe.Pointer {
	newMyData := new(MyData)
	newMyData.Size = size
	newMyData.Arr = make([]uint32, size)
	fmt.Print("Input: ")
	for i := 0; i < int(size); i++ {
		newMyData.Arr[i] = 10 * uint32(i+1)
		fmt.Print(newMyData.Arr[i], " ")
	}
	fmt.Println()
	return unsafe.Pointer(newMyData)
}

/* Function that serializes an instance of MyData */
func Serialize(buf unsafe.Pointer) (unsafe.Pointer, uint32) {
	mydata := (*MyData)(buf)
	serialBuf := make([]uint32, mydata.Size+1)
	serialBuf[0] = uint32(mydata.Size)

	var i uint32
	for i = 0; i < mydata.Size; i++ {
		serialBuf[i+1] = mydata.Arr[i]
	}

	retBuf := unsafe.Pointer(&serialBuf[0])
	bytes := (mydata.Size + 1) * 4

	return retBuf, bytes

}

/* Function that constructs an instance of MyData out of serialized data */
func Deserialize(buf unsafe.Pointer) unsafe.Pointer {
	sizeExtr := *((*uint32)(buf))

	/* Convert unsafe.Pointer to Slice */
	var slice []uint32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Data = uintptr(buf)
	header.Len = int(sizeExtr + 1)
	header.Cap = int(sizeExtr + 1)

	/* Reconstruct the structure */
	mydatabuf := new(MyData)
	mydatabuf.Size = sizeExtr
	mydatabuf.Arr = make([]uint32, sizeExtr)

	var i uint32
	for i = 0; i < sizeExtr; i++ {
		mydatabuf.Arr[i] = slice[i+1]
	}

	return unsafe.Pointer(mydatabuf)

}

func main() {
	/* Read User Args */
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <shared-object>\n", os.Args[0])
		return
	}

	path := os.Args[1]

	var session vaccel.Session
	var sharedObject vaccel.Resource

	err := sharedObject.Init(path, vaccel.ResourceLib)
	if err != vaccel.OK {
		fmt.Println("error creating shared object")
		os.Exit(int(err))
	}

	err = session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(int(err))
	}

	err = session.Register(&sharedObject)
	if err != vaccel.OK {
		fmt.Println("error registering resource with session")
		os.Exit(int(err))
	}

	/* Create the arg-lists */
	read := vaccel.ArgsInit(1)
	write := vaccel.ArgsInit(1)

	if read == nil || write == nil {
		fmt.Println("Error Creating the arg-lists")
		os.Exit(0)
	}

	/* Add a non-serialized arg */
	/* 10 20 30 40 50 */
	var numEntries uint32 = 5
	myDataPtr := NewMyData(numEntries)

	if read.AddNonSerialArg(myDataPtr, 0, Serialize) != vaccel.OK {
		fmt.Println("Error Adding Non-Serialized arg")
		os.Exit(0)
	}

	/* Define an expected argument */
	var uint32Size uint32 = 4
	expectedSize := int((numEntries + 1) * uint32Size)
	if write.ExpectNonSerialArg(expectedSize) != vaccel.OK {
		fmt.Println("Error defining expected arg")
		os.Exit(0)
	}

	/* Run the operation */
	err = vaccel.ExecWithResource(&session, &sharedObject, "mytestfunc_nonser", read, write)
	if err != vaccel.OK {
		fmt.Println("An error occurred while running the operation")
		os.Exit(err)
	}

	/* Extract the Output */
	outbuf := write.ExtractNonSerialArg(0, Deserialize)
	mydataOut := (*MyData)(outbuf)

	fmt.Print("Output: ")
	for i := 0; i < int(mydataOut.Size); i++ {
		fmt.Print(mydataOut.Arr[i], " ")
	}
	fmt.Println()

	/* Delete the lists */
	if write.Delete() != vaccel.OK || read.Delete() != vaccel.OK {
		fmt.Println("An error occurred in deletion of the arg-lists")
		os.Exit(0)
	}

	if session.Unregister(&sharedObject) != vaccel.OK {
		fmt.Println("An error occurred while unregistering the resource")
	}

	if sharedObject.Release() != vaccel.OK {
		fmt.Println("An error occurred while releasing the resource")
	}

	if session.Release() != vaccel.OK {
		fmt.Println("An error occurred while freeing the session")
	}
}
