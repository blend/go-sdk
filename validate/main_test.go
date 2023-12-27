/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

func none() error	{ return nil }

func some(err error) func() error	{ return func() error { return err } }
