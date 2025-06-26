// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {
	var session vaccel.Session

	err := session.Init(0)
	if err != vaccel.OK {
		fmt.Println("error initializing session")
		os.Exit(int(err))
	}

	err = vaccel.NoOp(&session)
	if err != vaccel.OK {
		fmt.Println("An error occurred while running the operation")
		os.Exit(err)
	}

	err = session.Release()
	if err != vaccel.OK {
		fmt.Println("An error occurred while freeing the session")
		os.Exit(err)
	}

}
