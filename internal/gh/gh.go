// Package gh is a small GitHub REST client, limited to the two endpoints the
// version checker needs: the tags of a repository and its releases. It is
// deliberately dependency free; a token is optional but raises the rate limit
// from 60 to 5000 requests an hour, which matters when checking a whole repository.
package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client talks to the GitHub API.
type Client struct {
	// BaseURL defaults to https://api.github.com.
	BaseURL string
	// Token is an optional personal access token, sent as a bearer token.
	Token string
	// UserAgent identifies the caller, as required by the GitHub API.
	UserAgent string
	// HTTPClient defaults to one with a 30s timeout.
	HTTPClient *http.Client
}

// New returns a client usable as is.
func New(token string) *Client {
	return &Client{
		BaseURL:    "https://api.github.com",
		Token:      token,
		UserAgent:  "binpot-versions",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Release is the subset of a GitHub release this package cares about.
type Release struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
}

type tag struct {
	Name string `json:"name"`
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	url := strings.TrimSuffix(c.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("building request for %s: %w", path, err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("getting %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("getting %s: %s (rate limited? set GH_TOKEN): %s",
				path, resp.Status, strings.TrimSpace(string(body)))
		}
		return fmt.Errorf("getting %s: %s", path, resp.Status)
	}
	// The body is decoded as a stream: a page of one hundred releases of a busy
	// repository such as helm/helm is well over a megabyte, and buffering with a
	// cap truncates the JSON.
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding %s: %w", path, err)
	}
	return nil
}

// Tags returns the tag names of a repository, newest first, reading at most
// maxPages pages of one hundred tags. GitHub orders tags by descending creation
// date, so a handful of pages covers the recent releases of even very busy
// repositories such as kubernetes/kubernetes.
func (c *Client) Tags(ctx context.Context, repo string, maxPages int) ([]string, error) {
	if maxPages < 1 {
		maxPages = 1
	}
	var names []string
	for page := 1; page <= maxPages; page++ {
		var tags []tag
		path := fmt.Sprintf("/repos/%s/tags?per_page=100&page=%d", repo, page)
		if err := c.get(ctx, path, &tags); err != nil {
			return names, err
		}
		for _, t := range tags {
			names = append(names, t.Name)
		}
		if len(tags) < 100 {
			break
		}
	}
	return names, nil
}

// Releases returns the most recent releases of a repository, drafts and
// pre-releases included, up to maxPages pages.
func (c *Client) Releases(ctx context.Context, repo string, maxPages int) ([]Release, error) {
	if maxPages < 1 {
		maxPages = 1
	}
	var out []Release
	for page := 1; page <= maxPages; page++ {
		var releases []Release
		path := fmt.Sprintf("/repos/%s/releases?per_page=100&page=%d", repo, page)
		if err := c.get(ctx, path, &releases); err != nil {
			return out, err
		}
		out = append(out, releases...)
		if len(releases) < 100 {
			break
		}
	}
	return out, nil
}
