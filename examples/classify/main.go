// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./classify <filename>")
		return
	}

	path := os.Args[1]

	var session vaccel.Session

	err := session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(err)
	}

	var outText string
	outText, err = vaccel.ImageClassificationFromFile(&session, path)
	if err != vaccel.OK {
		fmt.Println("Image Classification failed")
		os.Exit(err)
	}

	fmt.Println("Output(1): ", outText)

	imageBytes, e := os.ReadFile(path)
	if e != nil {
		fmt.Printf("Error reading file: %s\n", e)
		os.Exit(1)
	}

	outText, err = vaccel.ImageClassification(&session, imageBytes)
	if err != vaccel.OK {
		fmt.Println("Image Classification failed")
		os.Exit(err)
	}

	fmt.Println("Output(2): ", outText)

	err = session.Release()
	if err != vaccel.OK {
		fmt.Println("An error occurred while freeing the session")
		os.Exit(1)
	}
}
