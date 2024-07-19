/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptJSONBody sets the post body on the request.
func OptJSONBody(obj interface{}) Option {
	return RequestOption(webutil.OptJSONBody(obj))
}
