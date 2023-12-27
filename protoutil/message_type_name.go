/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package protoutil

import "google.golang.org/protobuf/proto"

// MessageTypeName returns the message type name for a given message.
func MessageTypeName(msg proto.Message) string {
	return string(msg.ProtoReflect().Type().Descriptor().FullName())
}
