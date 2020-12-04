/*
Package r2 is a rewrite of the sdk http request package that eschews fluent apis in favor of the options pattern.

The request returned by `r2.New()` i.e. `*r2.Request` holds everything required to send the request, including the http client reference, and a transport reference. If neither are specified, defaults are used (http.DefaultClient for the client, etc.)

To send a request, simply:

	resp, err := r2.New("http://example.com/").Do()
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...

You can specify additional options as a variadic list of `Opt...` functions:

	resp, err := r2.New("http://example.com",
		r2.OptPost(),
		r2.OptHeaderValue("X-Foo", "bailey"),
	).Do()

There are convenience methods on the request type that help with things like decoding types as json:

	meta, err := r2.New("http://example.com",
		r2.OptPost(),
		r2.OptHeaderValue("X-Foo", "bailey"),
	).JSON(&myObj)

Note that in the above, the JSON method also returns a closed version of the response for metadata purposes.

You can also fire and forget the request with the `.Discard()` method:

	meta, err := r2.New("http://example.com",
		r2.OptPost(),
		r2.OptHeaderValue("X-Foo", "bailey"),
	).Discard()


*/
package r2 // import "github.com/blend/go-sdk/r2"
