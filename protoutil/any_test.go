package protoutil

import (
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/protoutil/testdata"
	"github.com/blend/go-sdk/uuid"
)

func Test_Any(t *testing.T) {
	its := assert.New(t)

	original := newTestMessage()
	packed, err := Any(original)
	its.Nil(err)
	its.NotNil(packed)
	its.Equal(TypeURLPrefix+"testdata.Message", packed.TypeUrl)

	unpacked, err := FromAny(packed)
	its.Nil(err)
	its.NotNil(unpacked)
	typed, ok := unpacked.(*testdata.Message)
	its.True(ok)
	its.Equal(original.Uid, typed.Uid)
	its.Equal(FromTimestamp(original.TimestampUtc), FromTimestamp(typed.TimestampUtc))
	its.Equal(FromDuration(original.Elapsed), FromDuration(typed.Elapsed))
	its.Equal(original.StatusCode, typed.StatusCode)
	its.Equal(original.Value, typed.Value)
	its.Equal(original.Error, typed.Error)

	// from any handles nil ...
	unpacked, err = FromAny(nil)
	its.Equal(ex.Class("cannot unpack message from nil *any.Any"), ex.ErrClass(err))
	its.Nil(unpacked)

	// any handles bogus inputs
	var bogus proto.Message
	packed, err = Any(bogus)
	its.NotNil(err)
	its.Nil(packed)
}

func Test_FromTypeURL(t *testing.T) {
	its := assert.New(t)

	bareMessage, err := FromTypeURL("testdata.Message")
	its.Nil(err)
	its.NotNil(bareMessage)
	typed, ok := bareMessage.(*testdata.Message)
	its.True(ok)
	its.NotNil(typed)

	notFound, err := FromTypeURL(uuid.V4().String())
	its.NotNil(err)
	its.Nil(notFound)
}

func Test_TypeURL(t *testing.T) {
	its := assert.New(t)

	its.Equal(TypeURLPrefix+"testdata.Message", TypeURL(new(testdata.Message)))

	its.Equal("testdata.Message", MessageTypeFromTypeURL(TypeURL(new(testdata.Message))))
}
