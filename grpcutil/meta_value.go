/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import "google.golang.org/grpc/metadata"

// MetaValue returns a value from a metadata set.
func MetaValue(md metadata.MD, key string) string {
	if values, ok := md[key]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}
