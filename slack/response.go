/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
