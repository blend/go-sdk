package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/blend/go-sdk/breaker"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

func main() {

	mockServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if rand.Float64() > 0.5 {
			http.Error(rw, "should fail", http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!\n")
		return
	}))
	defer mockServer.Close()

	callServer := func(ctx context.Context) error {
		res, err := r2.New(mockServer.URL).Do()
		if err != nil {
			return err
		}
		if res.StatusCode >= 300 {
			return fmt.Errorf("non 200 status code returned from remote")
		}
		return nil
	}

	cb := breaker.MustNew(callServer,
		breaker.OptOpenExpiryInterval(5*time.Second),
	)

	var err error
	for x := 0; x < 1024; x++ {
		if err = cb.Execute(context.Background()); err != nil {
			fmt.Printf("circuit breaker error: %v\n", err)
		} else {
			fmt.Printf("circuit breaker call ok\n")
		}
		if ex.Is(err, breaker.ErrOpenState) {
			time.Sleep(5 * time.Second)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
