/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GithubClient is a client to make requests to github via the api.
type GithubClient struct {
	Addr	string
	Token	string
}

// CreateURL creates a fully qualified url with a given path.
func (ghc GithubClient) CreateURL(path string) string {
	parsedAddr, _ := url.Parse(ghc.Addr)
	parsedAddr.Path = path
	return parsedAddr.String()
}

// Do performs a request.
func (ghc GithubClient) Do(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", ghc.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if statusCode := res.StatusCode; statusCode < http.StatusOK || statusCode > 299 {
		defer func() { _ = res.Body.Close() }()
		contents, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("non-200 returned from remote; %s", string(contents))
	}
	return res, nil
}

// Post makes a post request.
func (ghc GithubClient) Post(ctx context.Context, path string, input, output interface{}) error {
	inputContents, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", ghc.CreateURL(path), bytes.NewReader(inputContents))
	if err != nil {
		return err
	}
	res, err := ghc.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := json.NewDecoder(res.Body).Decode(output); err != nil {
		return err
	}
	return nil
}

// Get makes a get request.
func (ghc GithubClient) Get(ctx context.Context, path string, output interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", ghc.CreateURL(path), nil)
	if err != nil {
		return err
	}
	res, err := ghc.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()
	if err := json.NewDecoder(res.Body).Decode(output); err != nil {
		return err
	}
	return nil
}

// UserExists fetches a user and returns if that user exists or not.
func (ghc GithubClient) UserExists(ctx context.Context, username string) error {
	output := make(map[string]interface{})
	username = strings.TrimPrefix(strings.TrimSpace(username), "@")
	return ghc.Get(ctx, fmt.Sprintf("/api/v3/users/%s", username), &output)
}

// TeamExists fetches a team and returns if that team exists or not.
func (ghc GithubClient) TeamExists(ctx context.Context, teamName string) error {
	output := make(map[string]interface{})
	parts := strings.Split(teamName, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid team name: %q", teamName)
	}
	org := strings.TrimPrefix(strings.TrimSpace(parts[0]), "@")
	team := strings.TrimSpace(parts[1])
	return ghc.Get(ctx, fmt.Sprintf("/api/v3/orgs/%s/teams/%s", org, team), &output)
}
