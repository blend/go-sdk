package certutil

import (
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestCreateSelfServerCert(t *testing.T) {
	assert := assert.New(t)

	notAfter := time.Date(2016, 02, 03, 12, 0, 0, 0, time.UTC)
	cert, err := CreateSelfServerCert("foo.bar.com", OptSubjectOrganization("the goods"), OptNotAfter(notAfter))
	assert.Nil(err)
	assert.NotNil(cert)
	assert.Equal([]string{"the goods"}, cert.Certificates[0].Subject.Organization)
	assert.Equal("foo.bar.com", cert.Certificates[0].Subject.CommonName)
	assert.Equal(notAfter, cert.Certificates[0].NotAfter)
}
