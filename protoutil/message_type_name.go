/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import "google.golang.org/protobuf/proto"

// MessageTypeName returns the message type name for a given message.
func MessageTypeName(msg proto.Message) string {
	return string(msg.ProtoReflect().Type().Descriptor().FullName())
}
