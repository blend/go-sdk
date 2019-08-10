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

func createCaller(results chan Result, url string, opts ...r2.Option) breaker.ActionProvider {
	return func() breaker.Action {
		return func(ctx context.Context) error {
			res, err := r2.New(url, opts...).Do()
			if err != nil {
				return err
			}
			defer res.Body.Close()
			if res.StatusCode >= 300 {
				return fmt.Errorf("non 200 status code returned from remote")
			}
			var result Result
			json.NewDecoder(res.Body).Decode(&result)
			results <- result
			return nil
		}
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

	results := make(chan Result, 1)
	cb := breaker.MustNew(createCaller(results, mockServer.URL),
		breaker.OptOpenExpiryInterval(5*time.Second),
	)

	var err error
	for x := 0; x < 1024; x++ {
		if err = cb.Execute(context.Background()); err != nil {
			fmt.Printf("circuit breaker error: %v\n", err)
			if ex.Is(err, breaker.ErrOpenState) {
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			fmt.Printf("result: %v\n", <-results)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
