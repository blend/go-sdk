/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_StatementInterceptorChain(t *testing.T) {
	its := assert.New(t)

	var calls []string
	a := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "a")
		return statement + "a", nil
	}

	b := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "b")
		return statement + "b", nil
	}

	c := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "c")
		return statement + "c", nil
	}

	chain := StatementInterceptorChain(a, b, c)
	statement, err := chain(context.TODO(), "foo", "bar")
	its.Nil(err)
	its.Equal("barabc", statement)
	its.Equal([]string{"a", "b", "c"}, calls)
}

func Test_StatementInterceptorChain_Errors(t *testing.T) {
	its := assert.New(t)

	var calls []string
	a := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "a")
		return statement + "a", nil
	}

	b := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "b")
		return statement + "b", fmt.Errorf("this is just a test")
	}

	c := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "c")
		return statement + "c", nil
	}

	chain := StatementInterceptorChain(a, b, c)
	statement, err := chain(context.TODO(), "foo", "bar")
	its.NotNil(err)
	its.Equal("barab", statement)
	its.Equal([]string{"a", "b"}, calls)
}

func Test_StatementInterceptorChain_Empty(t *testing.T) {
	its := assert.New(t)

	chain := StatementInterceptorChain()
	statement, err := chain(context.TODO(), "foo", "bar")
	its.Nil(err)
	its.Equal("bar", statement)
}

func Test_StatementInterceptorChain_Single(t *testing.T) {
	its := assert.New(t)

	var calls []string
	a := func(_ context.Context, label, statement string) (string, error) {
		calls = append(calls, "a")
		return statement + "a", nil
	}

	chain := StatementInterceptorChain(a)
	statement, err := chain(context.TODO(), "foo", "bar")
	its.Nil(err)
	its.Equal("bara", statement)
	its.Equal([]string{"a"}, calls)
}
