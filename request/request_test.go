package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

type statusObject struct {
	Status string `json:"status" xml:"status"`
}

func statusOkObject() statusObject {
	return statusObject{"ok!"}
}

type testObject struct {
	ID           int       `json:"id" xml:"id"`
	Name         string    `json:"name" xml:"name"`
	TimestampUtc time.Time `json:"timestamp_utc" xml:"timestamp_utc"`
	Value        float64   `json:"value" xml:"value"`
}

type errorObject struct {
	Code    int    `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
}

func newTestObject() testObject {
	to := testObject{}
	to.ID = rand.Int()
	to.Name = fmt.Sprintf("Test Object %d", to.ID)
	to.TimestampUtc = time.Now().UTC()
	to.Value = rand.Float64()
	return to
}

func newErrorObject() errorObject {
	err := errorObject{}
	err.Code = 1
	err.Message = "error message"
	return err
}

func okMeta() *ResponseMeta {
	return &ResponseMeta{StatusCode: http.StatusOK}
}

func noContentMeta() *ResponseMeta {
	return &ResponseMeta{StatusCode: http.StatusNoContent}
}

func errorMeta() *ResponseMeta {
	return &ResponseMeta{StatusCode: http.StatusInternalServerError}
}

func notFoundMeta() *ResponseMeta {
	return &ResponseMeta{StatusCode: http.StatusNotFound}
}

func writeHeader(w http.ResponseWriter, meta *ResponseMeta) {
	if len(meta.ContentType) > 0 {
		w.Header().Set(HeaderContentType, meta.ContentType)
	} else {
		w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
	}
	for key, value := range meta.Headers {
		w.Header().Set(key, strings.Join(value, ";"))
	}
	w.WriteHeader(meta.StatusCode)
}

func writeJSON(w http.ResponseWriter, meta *ResponseMeta, response interface{}) error {
	bytes, err := json.Marshal(response)
	if err == nil {
		writeHeader(w, meta)

		count, err := w.Write(bytes)
		if count == 0 {
			return exception.New("writeJSON didnt write any bytes")
		}
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func mockEchoEndpoint(meta *ResponseMeta) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		writeHeader(w, meta)

		defer r.Body.Close()
		bytes, _ := ioutil.ReadAll(r.Body)
		w.Write(bytes)
	})
}

type validationFunc func(r *http.Request)

func mockEndpoint(meta *ResponseMeta, returnWithObject interface{}, validations validationFunc) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		if validations != nil {
			validations(r)
		}

		writeJSON(w, meta, returnWithObject)
	})
}

func mockTLSEndpoint(meta *ResponseMeta, returnWithObject interface{}, validations validationFunc) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		if validations != nil {
			validations(r)
		}

		writeJSON(w, meta, returnWithObject)
	})
}

func mockNoContentEndpoint(meta *ResponseMeta, validations validationFunc) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		if validations != nil {
			validations(r)
		}

		writeHeader(w, meta)
	})
}

func getMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func getTLSMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewTLSServer(handler)
}

func TestCreateHttpRequestWithUrl(t *testing.T) {
	assert := assert.New(t)
	sr := New().MustWithRawURL("http://localhost:5001/api/v1/path/2?env=dev&foo=bar")

	assert.Equal("http", sr.scheme)
	assert.Equal("localhost:5001", sr.host)
	assert.Equal("GET", sr.method)
	assert.Equal("/api/v1/path/2", sr.path)
	assert.Equal(2, len(sr.query))
	assert.Equal([]string{"dev"}, sr.query["env"])
	assert.Equal([]string{"bar"}, sr.query["foo"])
}

func TestHttpGet(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, nil)
	testObject := testObject{}
	meta, err := New().AsGet().MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestHttpGetSetsRequestStartTime(t *testing.T) {
	assert := assert.New(t)
	testStart := time.Now().UTC()
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, nil)
	testObject := testObject{}
	req := New().AsGet().MustWithRawURL(ts.URL)
	meta, err := req.JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
	assert.True(req.Meta().StartTime.Equal(testStart) || req.Meta().StartTime.After(testStart))
}

func TestHttpGetWithErrorHandler(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newErrorObject()
	ts := mockEndpoint(errorMeta(), returnedObject, nil)
	testObject := testObject{}
	errorObject := errorObject{}
	meta, err := New().AsGet().MustWithRawURL(ts.URL).JSONWithErrorHandler(&testObject, &errorObject)
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode)
	assert.Equal(returnedObject, errorObject)
}

func TestHttpGetWithFullMockedError(t *testing.T) {
	assert := assert.New(t)
	url := "http://localhost:5001/api/v2/err"
	req := New().WithMethod("GET").MustWithRawURL(url).WithMockProvider(
		func(_ *Request) *MockedResponse {
			return &MockedResponse{
				Meta: ResponseMeta{StatusCode: http.StatusInternalServerError},
				Err:  errors.New("this is a test error"),
			}
		})
	testObject := testObject{}
	errorObject := errorObject{}
	meta, err := req.JSONWithErrorHandler(&testObject, &errorObject)
	assert.NotNil(err)
	assert.Equal("this is a test error", err.Error())
	assert.Equal(http.StatusInternalServerError, meta.StatusCode)
}

func TestHttpGetWithExpiringTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("This test involves a 500ms timeout.")
	}

	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, func(r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
	})
	testObject := testObject{}

	before := time.Now()
	_, err := New().AsGet().WithTimeout(250 * time.Millisecond).MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	after := time.Now()

	diff := after.Sub(before)
	assert.NotNil(err)
	assert.True(diff < 260*time.Millisecond, "Timeout was ineffective.")
}

func TestHttpGetWithTimeout(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, func(r *http.Request) {
		assert.Equal("GET", r.Method)
	})
	testObject := testObject{}
	meta, err := New().AsGet().WithTimeout(250 * time.Millisecond).MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestHttpGetNoContent(t *testing.T) {
	assert := assert.New(t)
	emptyObject := testObject{}
	ts := mockNoContentEndpoint(noContentMeta(), nil)
	testObject := testObject{}
	meta, err := New().AsGet().MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, meta.StatusCode)
	assert.Equal(emptyObject, testObject)
}

func TestHttpGetNoContentWithErrorHandler(t *testing.T) {
	assert := assert.New(t)
	emptyObject := testObject{}
	ts := mockNoContentEndpoint(noContentMeta(), nil)
	errorObject := testObject{}
	testObject := testObject{}
	meta, err := New().AsGet().MustWithRawURL(ts.URL).JSONWithErrorHandler(&testObject, &errorObject)
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, meta.StatusCode)
	assert.Equal(emptyObject, testObject)
	assert.Equal(emptyObject, errorObject)
}

func TestTlsHttpGet(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockTLSEndpoint(okMeta(), returnedObject, nil)
	testObject := testObject{}
	meta, err := New().AsGet().MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestHttpPostWithPostData(t *testing.T) {
	assert := assert.New(t)

	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, func(r *http.Request) {
		value := r.PostFormValue("foo")
		assert.Equal("bar", value)
	})

	testObject := testObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithPostData("foo", "bar").JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestHttpPostWithBasicAuth(t *testing.T) {
	assert := assert.New(t)

	ts := mockEndpoint(okMeta(), statusOkObject(), func(r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.True(ok)
		assert.Equal("test_user", username)
		assert.Equal("test_password", password)
	})

	testObject := statusObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithBasicAuth("test_user", "test_password").WithPostBody([]byte(`{"status":"ok!"}`)).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal("ok!", testObject.Status)
}

func TestHttpPostWithHeader(t *testing.T) {
	assert := assert.New(t)

	ts := mockEndpoint(okMeta(), statusOkObject(), func(r *http.Request) {
		value := r.Header.Get("test_header")
		assert.Equal(value, "foosballs")
	})

	testObject := statusObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithHeader("test_header", "foosballs").WithPostBody([]byte(`{"status":"ok!"}`)).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal("ok!", testObject.Status)
}

func TestHttpPostWithCookies(t *testing.T) {
	assert := assert.New(t)

	cookie := &http.Cookie{
		Name:     "test",
		Value:    "foosballs",
		Secure:   true,
		HttpOnly: true,
		Path:     "/test",
		Expires:  time.Now().UTC().AddDate(0, 0, 30),
	}

	ts := mockEndpoint(okMeta(), statusOkObject(), func(r *http.Request) {
		readCookie, readCookieErr := r.Cookie("test")
		assert.Nil(readCookieErr)
		assert.Equal(cookie.Value, readCookie.Value)
	})

	testObject := statusObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithCookie(cookie).WithPostBody([]byte(`{"status":"ok!"}`)).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal("ok!", testObject.Status)
}

func TestHttpPostWithJSONBody(t *testing.T) {
	assert := assert.New(t)

	returnedObject := newTestObject()
	ts := mockEchoEndpoint(okMeta())

	testObject := testObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithPostBodyAsJSON(&returnedObject).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestHttpPostWithXMLBody(t *testing.T) {
	assert := assert.New(t)

	returnedObject := newTestObject()
	ts := mockEchoEndpoint(okMeta())

	testObject := testObject{}
	meta, err := New().AsPost().MustWithRawURL(ts.URL).WithPostBodyAsXML(&returnedObject).XMLWithMeta(&testObject)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returnedObject, testObject)
}

func TestMockedRequests(t *testing.T) {
	assert := assert.New(t)

	ts := mockEndpoint(errorMeta(), nil, func(r *http.Request) {
		assert.True(false, "This shouldnt run in a mocked context.")
	})

	verifyString, meta, err := New().AsPut().WithPostBody([]byte("foobar")).MustWithRawURL(ts.URL).WithMockProvider(func(_ *Request) *MockedResponse {
		return &MockedResponse{Meta: *okMeta(), Res: []byte("ok!")}
	}).StringWithMeta()

	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal("ok!", verifyString)
}

func TestHandlers(t *testing.T) {
	assert := assert.New(t)

	ts := mockEchoEndpoint(okMeta())

	var calledRequest, calledResponse bool
	err := New().AsPut().WithPostBody([]byte("foobar")).MustWithRawURL(ts.URL).WithRequestHandler(func(req *Request) {
		calledRequest = true
	}).WithResponseHandler(func(req *Request, res *ResponseMeta, contents []byte) {
		calledResponse = true
	}).Execute()

	assert.Nil(err)
	assert.True(calledRequest)
	assert.True(calledResponse)
}

func TestRequestLogger(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, func(r *http.Request) {
		assert.Equal("GET", r.Method)
	})

	buffer := bytes.NewBuffer(nil)
	log := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(buffer).WithUseColor(false).WithShowTimestamp(false))
	defer log.Close()

	testObject := testObject{}
	_, err := New().WithLogger(log).AsGet().MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)

	log.Drain()
	assert.True(strings.HasPrefix(buffer.String(), "[request] GET http://127.0.0.1"), buffer.String())
}

func TestClientTrace(t *testing.T) {
	assert := assert.New(t)
	returnedObject := newTestObject()
	ts := mockEndpoint(okMeta(), returnedObject, func(r *http.Request) {
		assert.Equal("GET", r.Method)
	})

	receivedByte := false
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			receivedByte = true
		},
	}

	testObject := testObject{}
	_, err := New().WithClientTrace(trace).AsGet().MustWithRawURL(ts.URL).JSONWithMeta(&testObject)
	assert.Nil(err)
	assert.True(receivedByte)
}

func TestRequestEnforcesTransportRequirements(t *testing.T) {
	assert := assert.New(t)

	req := New().AsGet().WithHost("foo.com").WithTLSSkipVerify(true)
	_, _, err := req.BytesWithMeta()
	assert.True(exception.Is(err, ErrRequiresTransport))
	assert.Nil(req.Transport())
}

func TestRequestWithQueryString(t *testing.T) {
	assert := assert.New(t)

	req, err := New().AsGet().WithRawURL("http://foo.bar.com")
	assert.Nil(err)
	req = req.WithQueryString("foo", "bar")
	req = req.WithQueryString("buzz", "fuzz")

	full, err := req.Request()
	assert.Nil(err)
	assert.NotNil(full.URL)
	assert.NotNil(full.URL.Query())
	assert.Equal("bar", full.URL.Query().Get("foo"))
	assert.Equal("fuzz", full.URL.Query().Get("buzz"))
}

func TestResponseRequiresTransport(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(New().AsGet().MustWithRawURL("https://foo.bar.com").WithTLSSkipVerify(true).Execute())
}

func TestResponseAppliesTransportDefaults(t *testing.T) {
	assert := assert.New(t)

	ts := mockEndpoint(okMeta(), nil, func(r *http.Request) {
		assert.Equal("GET", r.Method)
	})

	xport := &http.Transport{}
	assert.Nil(New().AsGet().MustWithRawURL(ts.URL).WithTransport(xport).WithTLSSkipVerify(true).Execute())

	assert.NotNil(xport.TLSClientConfig)
	assert.True(xport.TLSClientConfig.InsecureSkipVerify)

	assert.Nil(New().AsGet().MustWithRawURL(ts.URL).WithTransport(xport).WithTLSSkipVerify(false).Execute())

	assert.NotNil(xport.TLSClientConfig)
	assert.False(xport.TLSClientConfig.InsecureSkipVerify)
}

var (
	_ Tracer = (*mockTracer)(nil)
)

type mockTracer struct {
	OnStart  func(*http.Request)
	OnFinish func(*http.Request, *ResponseMeta, error)
}

func (mt mockTracer) Start(req *http.Request) TraceFinisher {
	if mt.OnStart != nil {
		mt.OnStart(req)
	}
	return &mockTraceFinisher{parent: &mt}
}

type mockTraceFinisher struct {
	parent *mockTracer
}

func (mtf mockTraceFinisher) Finish(req *http.Request, meta *ResponseMeta, err error) {
	mtf.parent.OnFinish(req, meta, err)
}

type testKey int

const testKeyID testKey = iota

func TestRequestTracer(t *testing.T) {
	assert := assert.New(t)

	MockResponseFromString("GET", "https://test.com/foo", 200, "just a test")
	defer ClearMockedResponses()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var hasValue bool
	req := New().
		MustWithRawURL("https://test.com/foo").
		WithMockProvider(MockedResponseInjector).
		WithTracer(mockTracer{
			OnStart: func(r *http.Request) {
				defer wg.Done()
				(*r) = *r.WithContext(context.WithValue(r.Context(), testKeyID, "bar"))
			},
			OnFinish: func(r *http.Request, meta *ResponseMeta, err error) {
				defer wg.Done()
				hasValue = r.Context().Value(testKeyID) != nil
			},
		})

	assert.Nil(req.Execute())
	wg.Wait()
	assert.True(hasValue)
}

func TestRequestWithHeaders(t *testing.T) {
	assert := assert.New(t)

	req := New().WithHeader("foo", "bar").WithHeaders(http.Header{"buzz": []string{"wuzz"}})
	assert.Equal("bar", req.Header().Get("foo"))
	assert.Equal("wuzz", req.Header().Get("buzz"))
}

func TestRequestWithPostedFile(t *testing.T) {
	assert := assert.New(t)

	fileContents := bytes.NewBuffer([]byte(`this is only a test`))

	r := New().MustWithRawURL("http://foo.bar.com/hello").WithPostedFile("testFile", "testFile.txt", fileContents)

	assert.NotEmpty(r.postedFiles)

	postBody, err := r.PostBody()
	assert.Nil(err)
	assert.NotNil(postBody)

	contents, err := ioutil.ReadAll(postBody)
	assert.Nil(err)
	assert.NotEmpty(contents)
}

func TestRequestWithPostedFileIntegration(t *testing.T) {
	assert := assert.New(t)

	server := getMockServer(func(w http.ResponseWriter, req *http.Request) {
		// assert posted files exist on request
		file, _, err := req.FormFile("testFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if file == nil {
			http.Error(w, "file not found", http.StatusBadRequest)
			return
		}

		contents, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if string(contents) != `this is only a test` {
			http.Error(w, "wrong contents", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})
	defer server.Close()

	fileContents := bytes.NewBuffer([]byte(`this is only a test`))
	r := New().MustWithRawURL(server.URL).WithPostedFile("testFile", "testFile.txt", fileContents)
	meta, err := r.ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
}

func TestRequestWithClient(t *testing.T) {
	assert := assert.New(t)
	c := &http.Client{}
	r := New().WithClient(c)
	assert.Equal(c, r.client)

	c = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r = New().WithClient(c)
	assert.Equal(c, r.client)
}

func TestResponse(t *testing.T) {
	assert := assert.New(t)
	req := New()
	_, err := req.Response()
	assert.NotNil(err)
	assert.NotNil(req.client)
	assert.Equal(&http.Client{}, req.client)

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	meta := okMeta()
	meta.StatusCode = 302
	req = New().WithClient(c).WithMockProvider(func(req *Request) *MockedResponse {
		return &MockedResponse{
			Meta: *meta,
			Res:  []byte("ok!"),
		}
	})
	resp, err := req.Response()
	assert.Nil(err)
	assert.NotNil(resp)
}
