// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/nubificus/vaccel-go/vaccel"
)

func main() {

	/* Session */
	var session vaccel.Session
	err := vaccel.SessionInit(&session, 0)

	if err != 0 {
		fmt.Println("error initializing session")
		os.Exit(int(err))
	}

	/* Run the operation */
	err = vaccel.NoOp(&session)

	if err != 0 {
		fmt.Println("An error occurred while running the operation")
		os.Exit(err)
	}

	/* Free Session */
	if vaccel.SessionRelease(&session) != 0 {
		fmt.Println("An error occurred while freeing the session")
	}

}
