r2
==

R2 is an expressive and minimal library wrapping http clients that handles some common options (like timeouts, transports etc.) in a more unified fashion.

## Philosophy

Departing from "Fluent APIs", `go-sdk/r2` uses an "Options" based api for configuring and making http client requests.

Funadmentally, it means taking code that looked like:

```golang
res, err := request.New().
	WithVerb("POST").
	MustWithURL("https://www.google.com/robots.txt").
	WithHeaderValue("X-Authorization", "none").
	WithHeaderValue(request.HeaderContentType, request.ContentTypeApplicationJSON).
	WithBody([]byte(`{"status":"maybe?"}`)).
	Execute()
```

And refactors that to:

```golang
res, err := r2.New("https://www.google.com/robots.txt",
	r2.OptPost(),
	r2.OptHeaderValue("X-Authorization", "none").
	r2.OptHeaderValue(request.HeaderContentType, request.ContentTypeApplicationJSON),
	r2.OptBody([]byte(`{"status":"maybe?"}`)).Do()
```

The key difference here is making use of a variadic list of "Options" which are really just functions that satisfy the signature `func(*r2.Request) error`. This lets developers _extend_ the possible options that can be specified, vs. having a strictly hard coded list hung off the `request.Request` object, which require a PR to make changes to.

## Example

```golang
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func CustomOption() r2.Option {
	return func(r *r2.Request) {
		r.Client.Timeout = 10 * time.Millisecond
	}
}

func main() {
	_, err := r2.New("https://google.com",
		r2.Get(),
		r2.Timeout(500*time.Millisecond),
		r2.Header("X-Sent-By", "go-sdk/request2"),
		r2.CookieValue("ssid", "baileydog01"),
		CustomOption(),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```
