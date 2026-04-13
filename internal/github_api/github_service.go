package github_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IGitHubApi interface {
	RepoExists(ctx context.Context, repo string) (bool, error)
	GetLatestRelease(ctx context.Context, repo string) (string, error)
}

type GitHubApi struct {
	httpClient *http.Client
	token      string
}

func CreateGitHubApi(token string) IGitHubApi{
	return &GitHubApi{
		&http.Client{
			Timeout: 10 * time.Second,
		},
		token,
	}
}

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

func (c *GitHubApi) RepoExists(ctx context.Context, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	c.addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	case 403:
		return false, fmt.Errorf("rate limit exceeded or forbidden")
	default:
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}

func (c *GitHubApi) GetLatestRelease(ctx context.Context, repo string) (string, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/releases/latest",
		repo,
	)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	c.addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil // немає релізів
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github error: %d", resp.StatusCode)
	}

	var data releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.TagName, nil
}

func (c *GitHubApi) addHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-release-subscriber")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}