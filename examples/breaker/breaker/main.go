package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/breaker"
	"github.com/blend/go-sdk/r2"
)

// Result is a json thingy.
type Result struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createCaller(url string, opts ...r2.Option) breaker.Action {
	return func(ctx context.Context) (interface{}, error) {
		res, err := r2.New(url, opts...).Do()
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
	}
}

func main() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if rand.Float64() > 0.5 {
			http.Error(rw, "should fail", http.StatusInternalServerError)
			return
		}
		webutil.WriteJSON(rw, http.StatusOK, Result{1, "Foo"})
		return
	}))
	defer mockServer.Close()

	cb := breaker.MustNew(breaker.OptOpenExpiryInterval(5 * time.Second))
	var err error
	var res interface{}
	for x := 0; x < 1024; x++ {
		if res, err = cb.Do(context.Background(), createCaller(mockServer.URL)); err != nil {
			fmt.Printf("(%v) circuit breaker error: %v\n", cb.EvaluateState(context.Background()), err)
			if ex.Is(err, breaker.ErrOpenState) {
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			fmt.Printf("(%v) result: %v\n", cb.EvaluateState(context.Background()), res)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
