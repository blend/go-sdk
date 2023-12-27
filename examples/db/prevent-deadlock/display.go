/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"log"

	"github.com/blend/go-sdk/ex"
)

func displayError(err error) error {
	asMulti, ok := err.(*multiError)
	if !ok {
		return err
	}

	if asMulti == nil || len(asMulti.Errors) == 0 {
		return ex.New("Expected a non-nil / non-empty error")
	}

	log.Println("Error(s):")
	for _, err := range asMulti.Errors {
		log.Printf("- %#v\n", err)
	}

	return nil
}
