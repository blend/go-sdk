/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Interceptors(t *testing.T) {
	its := assert.New(t)

	var calls []string
	a := InterceptorFunc(func(i Actioner) Actioner {
		calls = append(calls, "a")
		return i
	})
	b := InterceptorFunc(func(i Actioner) Actioner {
		calls = append(calls, "b")
		return i
	})
	c := InterceptorFunc(func(i Actioner) Actioner {
		calls = append(calls, "c")
		return i
	})
	i := Interceptors(a, b, c)
	i.Intercept(new(NoopActioner))
	its.Equal([]string{"a", "b", "c"}, calls)
}

func Test_Interceptors_all_nil(t *testing.T) {
	its := assert.New(t)

	i := Interceptors(nil, nil, nil)
	its.Nil(i)
}

func Test_Interceptors_some_nil(t *testing.T) {
	its := assert.New(t)

	var calls []string
	a := InterceptorFunc(func(i Actioner) Actioner {
		calls = append(calls, "a")
		return i
	})
	c := InterceptorFunc(func(i Actioner) Actioner {
		calls = append(calls, "c")
		return i
	})
	i := Interceptors(nil, a, nil, c)
	i.Intercept(new(NoopActioner))
	its.Equal([]string{"a", "c"}, calls)
}
