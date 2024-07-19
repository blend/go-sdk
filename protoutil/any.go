/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import (
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/blend/go-sdk/ex"
)

// TypeURLPrefix is the type url prefix for type urls.
const TypeURLPrefix = "type.googleapis.com/"

// Any packs a message as an Any.
func Any(msg proto.Message) (*any.Any, error) {
	m, err := anypb.New(msg)
	if err != nil {
		return nil, ex.New(err)
	}
	return &any.Any{
		TypeUrl: m.TypeUrl,
		Value:   m.Value,
	}, nil
}

// FromAny unpacks a message any into a type message.
func FromAny(m *any.Any) (proto.Message, error) {
	if m == nil {
		return nil, ex.New("cannot unpack message from nil *any.Any")
	}
	return anypb.UnmarshalNew(&anypb.Any{
		TypeUrl: m.TypeUrl,
		Value:   m.Value,
	}, proto.UnmarshalOptions{
		AllowPartial: true,
	})
}

// FromTypeURL returns a message from the global types proto registry.
func FromTypeURL(typeURL string) (proto.Message, error) {
	if !strings.HasPrefix(typeURL, TypeURLPrefix) {
		typeURL = TypeURLPrefix + typeURL
	}
	mt, err := protoregistry.GlobalTypes.FindMessageByURL(typeURL)
	if err != nil {
		return nil, ex.New(err)
	}
	return mt.New().Interface(), nil
}

// TypeURL returns the typeURL for a given message.
//
// The bulk of this method was lifted from the anypb source.
func TypeURL(msg proto.Message) string {
	return TypeURLPrefix + string(MessageTypeName(msg))
}

// MessageTypeFromTypeURL returns the message type from a given any type url.
func MessageTypeFromTypeURL(typeURL string) string {
	return strings.TrimPrefix(typeURL, TypeURLPrefix)
}
