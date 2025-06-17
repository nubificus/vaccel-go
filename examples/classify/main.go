// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: %s <filename> [model]\n", os.Args[0])
		return
	}

	image := os.Args[1]
	imageBytes, e := os.ReadFile(image)
	if e != nil {
		fmt.Printf("Error reading file: %s\n", e)
		os.Exit(1)
	}

	var outText string
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

	outText, err = vaccel.ImageClassificationFromFile(&session, image)
	if err != vaccel.OK {
		fmt.Println("Image Classification failed")
		goto UnregisterResource
	}
	fmt.Println("Output(1): ", outText)

	outText, err = vaccel.ImageClassification(&session, imageBytes)
	if err != vaccel.OK {
		fmt.Println("Image Classification failed")
		goto UnregisterResource
	}
	fmt.Println("Output(2): ", outText)

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
