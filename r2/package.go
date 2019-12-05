/*
Package r2 is a rewrite of the request package that eschews fluent apis in favor of the options pattern.

To send a request, simply:

	resp, err := r2.New("http://example.com").Do()

Note: you must close the response body when finished with it:

	resp, err := r2.New("http://example.com/").Do()
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...

You can specify additional options as a variadic list of `Opt` functions:

	resp, err := r2.New("http://example.com",
		OptPost(),
		OptHeaderValue("X-Foo", "bailey"),
	).Do()

There are convenience methods on the request type that help with things like reading json:

	err := r2.New("http://example.com",
		OptPost(),
		OptHeaderValue("X-Foo", "bailey"),
	).JSON(&myObj)

The request object itself represents everything required to send the request, including
the http client reference, and a transport reference. If neither are specified, defaults
are used (http.DefaultClient for the client, etc.)

*/
package r2
