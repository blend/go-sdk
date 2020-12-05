package uuid

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/blend/go-sdk/assert"
)

func TestV4(t *testing.T) {
	m := make(map[string]bool)
	for x := 1; x < 32; x++ {
		uuid := V4()
		s := uuid.ToFullString()
		if m[s] {
			t.Errorf("NewRandom returned duplicated UUID %s\n", s)
		}
		m[s] = true
		if v := uuid.Version(); v != 4 {
			t.Errorf("Random UUID of version %v\n", v)
		}
	}
}

func makeTestUUIDv4(versionNumber byte, variant byte) UUID {
	return []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, versionNumber, 0x0, variant, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
}

func TestIsUUIDv4(t *testing.T) {
	assert := assert.New(t)

	valid := makeTestUUIDv4(0x40, 0x80)
	versionInvalid := makeTestUUIDv4(0xF0, 0x80)
	variantInvalid := makeTestUUIDv4(0x40, 0xF0)
	lengthInvalid := UUID([]byte{})

	assert.True(valid.IsV4())
	assert.False(variantInvalid.IsV4())
	assert.False(versionInvalid.IsV4())
	assert.False(lengthInvalid.IsV4())
}

func TestParseUUIDv4Valid(t *testing.T) {
	assert := assert.New(t)

	validShort := V4().ToShortString()
	validParsedShort, err := Parse(validShort)
	assert.Nil(err)
	assert.True(validParsedShort.IsV4())
	assert.Equal(validShort, validParsedShort.ToShortString())

	validFull := V4().ToFullString()
	validParsedFull, err := Parse(validFull)
	assert.Nil(err)
	assert.True(validParsedFull.IsV4())
	assert.Equal(validFull, validParsedFull.ToFullString())

	validBracedShort := fmt.Sprintf("{%s}", validShort)
	validParsedBracedShort, err := Parse(validBracedShort)
	assert.Nil(err)
	assert.True(validParsedBracedShort.IsV4())
	assert.Equal(validShort, validParsedBracedShort.ToShortString())

	validBracedFull := fmt.Sprintf("{%s}", validFull)
	validParsedBracedFull, err := Parse(validBracedFull)
	assert.Nil(err)
	assert.True(validParsedBracedFull.IsV4())
	assert.Equal(validFull, validParsedBracedFull.ToFullString())
}

func TestParseUUIDv4Invalid(t *testing.T) {
	assert := assert.New(t)

	/*
		_, err := Parse("")
		assert.NotNil(err, "should handle empty strings")
	*/

	_, err := Parse("fcae3946f75d+3258678bb5e6795a6d3")
	assert.NotNil(err, "should handle invalid characters")

	_, err = Parse("4f2e28b7b8f94b9eba1d90c4452")
	assert.NotNil(err, "should handle invalid length uuids")
}

type marshalTest struct {
	ID UUID `json:"id" yaml:"id"`
}

func TestJSONMarshalers(t *testing.T) {
	assert := assert.New(t)

	id := V4()
	rawJSON := []byte(fmt.Sprintf(`{"id":"%s"}`, id.ToFullString()))

	var testVal marshalTest
	assert.Nil(json.Unmarshal(rawJSON, &testVal))
	assert.Equal(id.String(), testVal.ID.String())

	newJSON, err := json.Marshal(testVal)
	assert.Nil(err)

	var verify marshalTest
	assert.Nil(json.Unmarshal(newJSON, &verify))
	assert.Equal(id.String(), verify.ID.String())
}

func TestJSONMarshalFormat(t *testing.T) {
	assert := assert.New(t)

	id := V4()

	marshalled, err := id.MarshalJSON()
	assert.Nil(err)
	assert.Equal(fmt.Sprintf("%q", id.ToFullString()), string(marshalled))
}

func TestYAMLMarshalers(t *testing.T) {
	assert := assert.New(t)

	id := V4()
	rawYAML := []byte(fmt.Sprintf(`id: "%s"`, id.ToFullString()))

	var testVal marshalTest
	assert.Nil(yaml.Unmarshal(rawYAML, &testVal))
	assert.Equal(id.String(), testVal.ID.String())

	newYAML, err := yaml.Marshal(testVal)
	assert.Nil(err)

	var verify marshalTest
	assert.Nil(yaml.Unmarshal(newYAML, &verify))
	assert.Equal(id.String(), verify.ID.String())
}
