/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package email

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMessageOption(t *testing.T) {
	assert := assert.New(t)

	m := Message{}

	assert.Empty(m.From)
	OptFrom("foo@bar.com")(&m)
	assert.Equal("foo@bar.com", m.From)

	assert.Empty(m.To)
	OptTo("buzz@bar.com", "fuzz@bar.com")(&m)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)

	assert.Empty(m.CC)
	OptCC("cc0@bar.com", "cc1@bar.com")(&m)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)

	assert.Empty(m.BCC)
	OptBCC("bcc0@bar.com", "bcc1@bar.com")(&m)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)

	assert.Empty(m.Subject)
	OptSubject("subject0")(&m)
	assert.Equal("subject0", m.Subject)

	assert.Empty(m.TextBody)
	OptTextBody("text body etc.")(&m)
	assert.Equal("text body etc.", m.TextBody)

	assert.Empty(m.HTMLBody)
	OptHTMLBody("html body etc.")(&m)
	assert.Equal("html body etc.", m.HTMLBody)

	assert.Equal("foo@bar.com", m.From)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)
	assert.Equal("text body etc.", m.TextBody)
	assert.Equal("html body etc.", m.HTMLBody)
}

func TestApplyMessageOptions(t *testing.T) {
	assert := assert.New(t)

	m := ApplyMessageOptions(Message{}, OptFrom("foo@bar.com"),
		OptTo("buzz@bar.com", "fuzz@bar.com"),
		OptCC("cc0@bar.com", "cc1@bar.com"),
		OptBCC("bcc0@bar.com", "bcc1@bar.com"),
		OptTextBody("text body etc."),
		OptHTMLBody("html body etc."))

	assert.Equal("foo@bar.com", m.From)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)
	assert.Equal("text body etc.", m.TextBody)
	assert.Equal("html body etc.", m.HTMLBody)
}
