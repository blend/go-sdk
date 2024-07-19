/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"

	"github.com/blend/go-sdk/ex"
)

func displayError(err error) error {
	asEx, ok := err.(*ex.Ex)
	if !ok {
		return err
	}

	if asEx == nil {
		return ex.New("Expected a non-nil error")
	}

	if asEx.Inner != nil {
		return ex.New("Did not expect an inner error")
	}

	log.Println("Error(s):")
	log.Printf("- Message: %q\n", asEx.Message)
	log.Printf("- %#v\n", asEx.Class)
	return nil
}
