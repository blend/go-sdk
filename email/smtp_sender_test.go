/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package email

import "github.com/blend/go-sdk/configutil"

var (
	_	configutil.Resolver	= (*SMTPSender)(nil)
	_	configutil.Resolver	= (*SMTPPlainAuth)(nil)
)
