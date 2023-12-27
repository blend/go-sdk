/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCheckKey(t *testing.T) {
	its := assert.New(t)

	its.Nil(CheckKey("foo"))
	its.Nil(CheckKey("bar/foo"))
	its.Nil(CheckKey("bar.io/foo"))
	its.NotNil(CheckKey("_foo"))
	its.NotNil(CheckKey("-foo"))
	its.NotNil(CheckKey("foo-"))
	its.NotNil(CheckKey("foo_"))
	its.NotNil(CheckKey("bar/foo/baz"))

	its.NotNil(CheckKey(""), "should error on empty keys")

	its.NotNil(CheckKey("/foo"), "should error on empty dns prefixes")
	superLongDNSPrefixed := fmt.Sprintf("%s/%s", strings.Repeat("a", MaxLabelKeyDNSSubdomainLen), strings.Repeat("a", MaxLabelKeyLen))
	its.Nil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxLabelKeyDNSSubdomainLen+1), strings.Repeat("a", MaxLabelKeyLen))
	its.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxLabelKeyDNSSubdomainLen+1), strings.Repeat("a", MaxLabelKeyLen+1))
	its.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxLabelKeyDNSSubdomainLen), strings.Repeat("a", MaxLabelKeyLen+1))
	its.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
}

func TestCheckKeyK8S(t *testing.T) {
	assert := assert.New(t)

	values := []string{
		"simple",
		"now-with-dashes",
		"1-starts-with-num",
		"1234",
		"simple/simple",
		"now-with-dashes/simple",
		"now-with-dashes/now-with-dashes",
		"now.with.dots/simple",
		"now-with.dashes-and.dots/simple",
		"1-num.2-num/3-num",
		"1234/5678",
		"1.2.3.4/5678",
		"Uppercase_Is_OK_123",
		"example.com/Uppercase_Is_OK_123",
		"requests.storage-foo",
		strings.Repeat("a", 63),
		strings.Repeat("a", 253) + "/" + strings.Repeat("b", 63),
	}
	badValues := []string{
		"nospecialchars%^=@",
		"cantendwithadash-",
		"-cantstartwithadash-",
		"only/one/slash",
		"example_com/abc",
		"example.com/",
		"Example.com/abc",
		"/simple",
		strings.Repeat("a", 64),
		strings.Repeat("a", 254) + "/abc",
	}
	for _, val := range values {
		assert.Nil(CheckKey(val), "input:", val)
	}
	for _, val := range badValues {
		assert.NotNil(CheckKey(val), "input:", val)
	}
}
