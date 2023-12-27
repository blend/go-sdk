/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/blend/go-sdk/ex"
)

// AWSAuth defines vault aws auth methods
type AWSAuth struct {
	CredentialProvider CredentialProvider
}

// NewAWSAuth creates a new AWS struct
func NewAWSAuth(opts ...AWSAuthOption) (*AWSAuth, error) {
	auth := &AWSAuth{
		CredentialProvider: GetIAMAuthCredentials,
	}
	var err error
	for _, opt := range opts {
		if err = opt(auth); err != nil {
			return nil, err
		}
	}
	return auth, nil
}

// AWSAuthOption mutates an AWSAuth instance
type AWSAuthOption func(*AWSAuth) error

// CredentialProvider defines the credential provider func interface
type CredentialProvider func(roleARN string) (*credentials.Credentials, error)

// OptAWSAuthCredentialProvider sets the credential provider
func OptAWSAuthCredentialProvider(cp CredentialProvider) AWSAuthOption {
	return func(a *AWSAuth) error {
		a.CredentialProvider = cp
		return nil
	}
}

// AWSIAMLogin returns a vault token given the instance role which invokes this function
func (a *AWSAuth) AWSIAMLogin(ctx context.Context, client HTTPClient, baseURL url.URL, roleName, roleARN, service, region string) (string, error) {
	stsRequest, err := a.GetCallerIdentitySignedRequest(roleARN, service, region)
	if err != nil {
		return "", ex.New(err)
	}

	request, err := createVaultLoginRequest(roleName, baseURL, stsRequest)
	if err != nil {
		return "", ex.New(err)
	}

	res, err := client.Do(request)
	if err != nil {
		return "", ex.New(err)
	}
	defer res.Body.Close()

	var response AWSAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", ex.New(err)
	}
	if len(response.Errors) > 0 {
		return "", ex.New("Error making aws get identity request", ex.OptMessagef("%+v", response.Errors))
	}

	return response.Auth.ClientToken, nil
}

// GetCallerIdentitySignedRequest gets a signed caller identity request
func (a *AWSAuth) GetCallerIdentitySignedRequest(roleARN, service, region string) (*http.Request, error) {
	credentials, err := a.CredentialProvider(roleARN)
	if err != nil {
		return nil, ex.New(err)
	}

	stsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials:	credentials,
			Region:		&region,
		},
	})
	if err != nil {
		return nil, ex.New(err)
	}

	svc := sts.New(stsSession)
	stsRequest, _ := svc.GetCallerIdentityRequest(nil)
	err = stsRequest.Sign()
	if err != nil {
		return nil, ex.New(err)
	}

	return stsRequest.HTTPRequest, nil
}

// GetIAMAuthCredentials is a credential provider to be passed in as input into the AWSAuth struct
func GetIAMAuthCredentials(roleARN string) (*credentials.Credentials, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, ex.New(err)
	}
	credentials := stscreds.NewCredentials(session, roleARN)
	return credentials, nil
}

func createVaultLoginRequest(roleName string, baseURL url.URL, request *http.Request) (*http.Request, error) {
	baseURL.Path = AWSAuthLoginPath
	stsHeaders, err := json.Marshal(request.Header)
	if err != nil {
		return nil, ex.New(err)
	}

	body := map[string]string{
		"role":				roleName,
		"iam_http_request_method":	MethodPost,
		"iam_request_url":		base64.StdEncoding.EncodeToString([]byte(request.URL.String())),
		"iam_request_body":		base64.StdEncoding.EncodeToString([]byte(STSGetIdentityBody)),
		"iam_request_headers":		base64.StdEncoding.EncodeToString(stsHeaders),
	}

	contents, err := json.Marshal(body)
	if err != nil {
		return nil, ex.New(err)
	}

	req := &http.Request{
		URL:	&baseURL,
		Method:	MethodPost,
		Body:	io.NopCloser(bytes.NewReader(contents)),
	}

	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(contents)
		return io.NopCloser(r), nil
	}

	req.ContentLength = int64(len(contents))
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set(HeaderContentType, ContentTypeApplicationJSON)

	return req, nil
}
