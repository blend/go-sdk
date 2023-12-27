/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slack

// PostMessageResponse is a slack response
type PostMessageResponse struct {
	OK		bool	`json:"ok"`
	Channel		string	`json:"channel,omitempty"`
	Timestamp	string	`json:"ts,omitempty"`
	Message		Message	`json:"message,omitempty"`
	Error		string	`json:"error,omitempty"`
}
