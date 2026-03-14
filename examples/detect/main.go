// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s <filename> [model]\n", os.Args[0])
		return
	}

	image := filepath.Clean(os.Args[1])
	imageBytes, e := os.ReadFile(image)
	if e != nil {
		fmt.Printf("Error reading file: %s\n", e)
		os.Exit(1)
	}

	var outImgName string
	var session vaccel.Session
	var resource vaccel.Resource

	err := session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(err)
	}

	if len(os.Args) == 3 {
		model := os.Args[2]
		err = resource.Init(model, vaccel.ResourceModel)
		if err != vaccel.OK {
			fmt.Printf("error initializing resource: %d\n", err)
			goto ReleaseSession
		}

		err = session.Register(&resource)
		if err != vaccel.OK {
			fmt.Printf("error while registering resource: %d\n", err)
			goto ReleaseResource
		}
	}

	outImgName, err = vaccel.ImageDetectionFromFile(&session, image)
	if err != vaccel.OK {
		fmt.Println("Image Detection failed")
		goto UnregisterResource
	}
	fmt.Println("Output(1): ", outImgName)

	outImgName, err = vaccel.ImageDetection(&session, imageBytes)
	if err != vaccel.OK {
		fmt.Println("Image Detection failed")
		goto UnregisterResource
	}
	fmt.Println("Output(2): ", outImgName)

UnregisterResource:
	if len(os.Args) == 3 {
		err = session.Unregister(&resource)
		if err != vaccel.OK {
			fmt.Printf("Could not unregister resource: %d\n", err)
		}
	}
ReleaseResource:
	if len(os.Args) == 3 {
		err = resource.Release()
		if err != vaccel.OK {
			fmt.Printf("Could not release resource: %d\n", err)
		}
	}
ReleaseSession:
	err = session.Release()
	if err != vaccel.OK {
		fmt.Println("An error occurred while freeing the session")
		os.Exit(1)
	}
}
