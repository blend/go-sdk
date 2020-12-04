package vault

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func mustURLf(format string, args ...interface{}) *url.URL {
	return webutil.MustParseURL(fmt.Sprintf(format, args...))
}

func TestVaultClientBackendKV(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	mountMetaJSON := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"2"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().WithString("GET", mustURLf("%s/v1/sys/internal/ui/mounts/secret/foo/bar", client.Remote.String()), mountMetaJSON)
	client.Client = m

	backend, err := client.backendKV(todo, "foo/bar")
	assert.Nil(err)
	assert.NotNil(backend)
}

func TestVaultClientGetVersion(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	mountMetaJSONV1 := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"1"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`
	mountMetaJSONV2 := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"2"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().
		WithString("GET", mustURLf("%s/v1/sys/internal/ui/mounts/secret/foo/bar", client.Remote.String()), mountMetaJSONV1)

	client.Client = m

	version, err := client.getVersion(todo, "foo/bar")
	assert.Nil(err)
	assert.Equal(Version1, version)

	m.WithString("GET", mustURLf("%s/v1/sys/internal/ui/mounts/secret/foo/bar", client.Remote.String()), mountMetaJSONV2)

	version, err = client.getVersion(todo, "foo/bar")
	assert.Nil(err)
	assert.Equal(Version2, version)
}

func TestVaultClientGetMountMeta(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	mountMetaJSON := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"2"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().WithString("GET", mustURLf("%s/v1/sys/internal/ui/mounts/secret/foo/bar", client.Remote.String()), mountMetaJSON)
	client.Client = m

	mountMeta, err := client.getMountMeta(todo, "secret/foo/bar")
	assert.Nil(err)
	assert.NotNil(mountMeta)
	assert.Equal(Version2, mountMeta.Data.Options["version"])
}

func TestVaultClientJSONBody(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	output, err := client.jsonBody(map[string]interface{}{
		"foo": "bar",
	})
	assert.Nil(err)
	defer output.Close()

	contents, err := ioutil.ReadAll(output)
	assert.Nil(err)
	assert.Equal("{\"foo\":\"bar\"}\n", string(contents))
}

func TestVaultClientReadJSON(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	jsonBody := bytes.NewBuffer([]byte(`{"foo":"bar"}`))

	output := map[string]interface{}{}
	assert.Nil(client.readJSON(jsonBody, &output))
	assert.Equal("bar", output["foo"])
}

func TestVaultClientCopyRemote(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	copy := client.copyRemote()
	copy.Host = "not_" + copy.Host

	anotherCopy := client.copyRemote()
	assert.NotEqual(anotherCopy.Host, copy.Host)
}

func TestVaultClientDiscard(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	assert.NotNil(client.discard(nil, fmt.Errorf("this is only a test")))

	assert.Nil(client.discard(client.jsonBody(map[string]interface{}{
		"foo": "bar",
	})))
}

func TestVaultCreateTransitKey(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	key := "key"

	m := NewMockHTTPClient().
		With(
			"POST",
			mustURLf("%s/v1/transit/keys/%s", client.Remote.String(), key),
			&http.Response{
				StatusCode: http.StatusNoContent,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
			},
		)
	client.Client = m

	err = client.CreateTransitKey(todo, "key")
	assert.Nil(err)
}

func TestVaultConfigureTransitKey(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	key := "key"

	m := NewMockHTTPClient().
		With(
			"POST",
			mustURLf("%s/v1/transit/keys/%s/config", client.Remote.String(), key),
			&http.Response{
				StatusCode: http.StatusNoContent,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
			},
		)
	client.Client = m

	err = client.ConfigureTransitKey(todo, "key", OptUpdateTransitDeletionAllowed(true))
	assert.Nil(err)
}

func TestVaultReadTransitKey(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	key := "key"
	keyMetaJSON := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"deletion_allowed":true,"exportable":false,"allow_plaintext_backup":false,"keys": {"1": 1442851412},"min_decryption_version": 1,"min_encryption_version": 0,"name": "foo"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().WithString("GET", mustURLf("%s/v1/transit/keys/%s", client.Remote.String(), key), keyMetaJSON)
	client.Client = m

	data, err := client.ReadTransitKey(todo, "key")
	assert.Nil(err)
	assert.Equal(true, data["deletion_allowed"])
}

func TestVaultDeleteTransitKey(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client, err := New()
	assert.Nil(err)

	key := "key"

	m := NewMockHTTPClient().
		With(
			"DELETE",
			mustURLf("%s/v1/transit/keys/%s", client.Remote.String(), key),
			&http.Response{
				StatusCode: http.StatusNoContent,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte{})),
			},
		)
	client.Client = m

	err = client.DeleteTransitKey(todo, "key")
	assert.Nil(err)
}

func TestVaultHandleRedirects(t *testing.T) {
	assert := assert.New(t)

	rawResponse := "{\"status\":\"ok!\"}\n"

	inner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(webutil.HeaderContentType, webutil.ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, rawResponse)
	}))
	defer inner.Close()
	outer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, inner.URL, http.StatusTemporaryRedirect)
	}))
	defer outer.Close()

	client, err := New(
		OptRemote(outer.URL),
	)
	assert.Nil(err)
	assert.NotNil(client)

	rawURL, err := url.Parse(outer.URL)
	assert.Nil(err)
	res, err := client.Client.Do(&http.Request{URL: rawURL})
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(http.StatusOK, res.StatusCode)

	contents, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)
	assert.Equal(rawResponse, string(contents))
}
