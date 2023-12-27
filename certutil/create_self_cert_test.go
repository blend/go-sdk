/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCreateSelfServerCert(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	notAfter := time.Date(2016, 02, 03, 12, 0, 0, 0, time.UTC)
	cert, err := CreateSelfServerCert("foo.bar.com", OptSubjectOrganization("the goods"), OptNotAfter(notAfter))
	assert.Nil(err)
	assert.NotNil(cert)
	assert.Equal([]string{"the goods"}, cert.Certificates[0].Subject.Organization)
	assert.Equal("foo.bar.com", cert.Certificates[0].Subject.CommonName)
	assert.Equal(notAfter, cert.Certificates[0].NotAfter)
}
