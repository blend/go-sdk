/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/blend/go-sdk/assert"
)

func TestAWSAuth_AWSIAMLogin(t *testing.T) {
	it := assert.New(t)

	client, err := New()
	it.Nil(err)

	sampleVaultResponse := `
		{
			"auth": {
				"renewable": true,
				"lease_duration": 1800000,
				"metadata": {
					"role_tag_max_ttl": "0",
					"instance_id": "i-de0f1344",
					"ami_id": "ami-fce36983",
					"role": "dev-role",
					"auth_type": "ec2"
				},
				"policies": ["default", "dev"],
				"accessor": "some-guid",
				"client_token": "my-test-token"
			}
		}`

	mockHTTPClient := NewMockHTTPClient().WithString("POST", mustURLf("%s/v1/auth/aws/login", client.Remote.String()), sampleVaultResponse)
	client.Client = mockHTTPClient
	authOpts := OptAWSAuthCredentialProvider(
		func(roleARN string) (*credentials.Credentials, error) {
			return credentials.NewStaticCredentials("id", "key", "session-token"), nil
		})
	client.AWSAuth, err = NewAWSAuth(authOpts)
	it.Nil(err)

	token, err := client.AWSAuth.AWSIAMLogin(context.TODO(), client.Client, *client.Remote, "roleName", "roleARN", "service", "us-east-1")
	it.Nil(err)
	it.Equal(token, "my-test-token")

}
