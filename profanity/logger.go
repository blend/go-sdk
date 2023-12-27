/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

// Logger are the methods required on the logger.
type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
}
