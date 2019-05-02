package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/uuid"
)

func main() {
	ca, _ := certutil.CreateCertificateAuthority()
	certBundle, _ := certutil.CreateServer(uuid.V4().String(), ca)

	certBundle.WriteCertPem(os.Stdout)
	fmt.Println()
	certBundle.WriteKeyPem(os.Stdout)
	fmt.Println()
}
