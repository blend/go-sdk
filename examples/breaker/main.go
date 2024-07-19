/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/blend/go-sdk/breaker"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

// Result is a json thingy.
type Result struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createUpstreamCaller(opts ...r2.Option) breaker.Actioner {
	return breaker.ActionerFunc(func(ctx context.Context, args interface{}) (interface{}, error) {
		res, err := r2.New(args.(string), opts...).Do()
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode >= 300 {
			return nil, fmt.Errorf("non 200 status code returned from remote")
		}
		var result Result
		json.NewDecoder(res.Body).Decode(&result)
		return result, nil
	})
}

var (
	flagNumCalls = flag.Int("num-calls", 1024, "The number of calls")
)

func init() {
	flag.Parse()
}

func main() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if rand.Float64() > 0.5 {
			http.Error(rw, "should fail", http.StatusInternalServerError)
			return
		}
		webutil.WriteJSON(rw, http.StatusOK, Result{1, "Foo"})
	}))
	defer mockServer.Close()

	b := breaker.New(
		breaker.OptOpenExpiryInterval(5 * time.Second),
	)
	cb := b.Intercept(createUpstreamCaller())

	var err error
	var res interface{}
	for x := 0; x < *flagNumCalls; x++ {
		if res, err = cb.Action(context.Background(), mockServer.URL); err != nil {
			fmt.Printf("(%v) circuit breaker error: %v\n", b.EvaluateState(context.Background()), err)
			if ex.Is(err, breaker.ErrOpenState) {
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			fmt.Printf("(%v) result: %v\n", b.EvaluateState(context.Background()), res)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
