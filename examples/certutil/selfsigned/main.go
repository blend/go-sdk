/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"
	"os"
	"time"

	"github.com/blend/go-sdk/certutil"
)

func main() {
	ca, err := certutil.CreateCertificateAuthority(
		certutil.OptSubjectOrganization("go-sdk"),
		certutil.OptNotAfter(time.Now().UTC().AddDate(0, 0, 30)),
	)
	if err != nil {
		log.Fatal(err)
	}

	server, err := certutil.CreateServer(
		"dev.local", ca,
		certutil.OptSubjectOrganization("go-sdk users"),
		certutil.OptNotAfter(time.Now().UTC().AddDate(0, 0, 15)),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = server.WriteCertPem(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = server.WriteKeyPem(os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
}
