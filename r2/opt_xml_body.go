/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptXMLBody sets the post body on the request.
func OptXMLBody(obj interface{}) Option {
	return RequestOption(webutil.OptXMLBody(obj))
}
