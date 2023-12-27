/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/uuid"
)

func main() {
	ca, _ := certutil.CreateCertificateAuthority(certutil.OptSubjectCommonName("go-sdk certificate authority"))

	ca.WriteCertPem(os.Stdout)
	fmt.Println()

	ca.WriteKeyPem(os.Stdout)
	fmt.Println()

	certBundle, _ := certutil.CreateServer(uuid.V4().String(), ca)

	certBundle.WriteCertPem(os.Stdout)
	fmt.Println()
	certBundle.WriteKeyPem(os.Stdout)
	fmt.Println()
}
