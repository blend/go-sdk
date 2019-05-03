package r2

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRequestNew(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.com/bar?buzz=fuzz")
	assert.NotNil(r)
	assert.Nil(r.Err)
	assert.Equal(MethodGet, r.Method)
	assert.NotNil(r.URL)
	assert.Equal("https://foo.com/bar?buzz=fuzz", r.URL.String())

	rErr := New("\n")
	assert.NotNil(rErr)
	assert.NotNil(rErr.Err)
}

func TestRequestDo(t *testing.T) {
	assert := assert.New(t)

	server := mockServerOK()
	defer server.Close()

	res, err := New(server.URL).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestDoAbortsOnError(t *testing.T) {
	assert := assert.New(t)

	var didCallServer bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		didCallServer = true
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	r := New(server.URL)
	r.Err = errors.New("this is only a test")
	_, err := r.Do()
	assert.NotNil(err)
	assert.Equal("this is only a test", err.Error())
	assert.False(didCallServer)
}

func TestRequestDoHeaders(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("foo"); value != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad value for foo: %#v\n", r.PostForm)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL, OptHeaderValue("foo", "bar")).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestDoQuery(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.URL.Query().Get("foo"); value != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad query value for foo: %#v\n", r.PostForm)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL, OptQueryValue("foo", "bar")).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestDoPostForm(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v!\n", err)
			return
		}
		if value := r.PostForm.Get("foo"); value != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad value for foo: %#v\n", r.PostForm)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer server.Close()

	res, err := New(server.URL,
		OptPost(),
		OptPostFormValue("foo", "bar"),
	).Do()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode, readString(res.Body))
}

func TestRequestDiscard(t *testing.T) {
	assert := assert.New(t)
	server := mockServerOK()
	defer server.Close()
	assert.Nil(New(server.URL).Discard())
}

func TestRequestDiscardWithResponse(t *testing.T) {
	assert := assert.New(t)
	server := mockServerOK()
	defer server.Close()
	res, err := New(server.URL).DiscardWithResponse()
	assert.Nil(err)
	assert.NotNil(res)
}

func TestRequestCopyTo(t *testing.T) {
	assert := assert.New(t)
	server := mockServerOK()
	defer server.Close()
	buf := new(bytes.Buffer)
	_, err := New(server.URL).CopyTo(buf)
	assert.Nil(err)
	assert.Equal("OK!\n", buf.String())
}

func TestRequestBytes(t *testing.T) {
	assert := assert.New(t)
	server := mockServerOK()
	defer server.Close()
	contents, err := New(server.URL).Bytes()
	assert.Nil(err)
	assert.Equal("OK!\n", contents)
}

func TestRequestBytesWithResponse(t *testing.T) {
	assert := assert.New(t)
	server := mockServerOK()
	defer server.Close()
	contents, meta, err := New(server.URL).BytesWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal("OK!\n", contents)
}

func TestRequestJSON(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\":\"ok!\"}\n")
	}))
	defer server.Close()

	var deserialized map[string]interface{}
	err := New(server.URL).JSON(&deserialized)
	assert.Nil(err)
	assert.Equal("ok!", deserialized["status"])
}

func TestRequestJSONWithResponse(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\":\"ok!\"}\n")
	}))
	defer server.Close()

	var deserialized map[string]interface{}
	res, err := New(server.URL).JSONWithResponse(&deserialized)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("ok!", deserialized["status"])
}

type xmlTestCase struct {
	Status string `xml:"status"`
}

func TestRequestXML(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		xml.NewEncoder(w).Encode(xmlTestCase{
			Status: "ok!",
		})
	}))
	defer server.Close()

	var deserialized xmlTestCase
	err := New(server.URL).XML(&deserialized)
	assert.Nil(err)
	assert.Equal("ok!", deserialized.Status)
}

func TestRequestXMLWithResponse(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		xml.NewEncoder(w).Encode(xmlTestCase{
			Status: "ok!",
		})
	}))
	defer server.Close()

	var deserialized xmlTestCase
	res, err := New(server.URL).XMLWithResponse(&deserialized)
	assert.Nil(err)
	assert.Equal("ok!", deserialized.Status)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestTracer(t *testing.T) {
	assert := assert.New(t)

	server := mockServerOK()
	defer server.Close()

	var didCallStart, didCallFinish bool
	tracer := MockTracer{
		StartHandler: func(_ *http.Request) {
			didCallStart = true
		},
		FinishHandler: func(_ *http.Request, _ *http.Response, _ time.Time, _ error) {
			didCallFinish = true
		},
	}
	assert.Nil(New(server.URL, OptTracer(tracer)).Discard())
	assert.True(didCallStart)
	assert.True(didCallFinish)
}

func TestRequestListeners(t *testing.T) {
	assert := assert.New(t)

	server := mockServerOK()
	defer server.Close()

	var didCallRequest1, didCallRequest2, didCallResponse1, didCallResponse2 bool
	assert.Nil(New(server.URL,
		OptOnRequest(func(_ *http.Request) error {
			didCallRequest1 = true
			return nil
		}),
		OptOnRequest(func(_ *http.Request) error {
			didCallRequest2 = true
			return nil
		}),
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error {
			didCallResponse1 = true
			return nil
		}),
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error {
			didCallResponse2 = true
			return nil
		}),
	).Discard())
	assert.True(didCallRequest1)
	assert.True(didCallRequest2)
	assert.True(didCallResponse1)
	assert.True(didCallResponse2)
}
