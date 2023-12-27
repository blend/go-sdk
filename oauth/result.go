/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package oauth

import "time"

// Result is the final result of the oauth exchange.
// It is the user profile of the user and the state information.
type Result struct {
	Response	Response
	Profile		Profile
	State		State
}

// Response is the response details from the oauth exchange.
type Response struct {
	AccessToken	string
	TokenType	string
	RefreshToken	string
	Expiry		time.Time
	HostedDomain	string
}
