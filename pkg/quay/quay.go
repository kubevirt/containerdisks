package quay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type QuayClient interface {
	Update(ctx context.Context, repository, description string, public bool) error
}

type quayClient struct {
	tokenFile string
	org       string
}

func (q *quayClient) base(repository string) url.URL {
	repoURL := url.URL{Host: "quay.io", Scheme: "https"}
	repoURL.Path = path.Join("/api/v1/repository", q.org, repository)
	return repoURL
}

func (q *quayClient) header() (http.Header, error) {
	rawToken, err := os.ReadFile(q.tokenFile)
	if err != nil {
		return http.Header{}, fmt.Errorf("error reading the quay token file: %v", err)
	}
	token := "Bearer " + strings.TrimSpace(string(rawToken))

	header := http.Header{}
	header.Add("Authorization", token)
	header.Add("Content-Type", "application/json")
	return header, nil
}

func (q *quayClient) json(ctx context.Context, method, repo, subresource string, jsonObj interface{}) error {
	content, err := json.Marshal(jsonObj)
	if err != nil {
		return fmt.Errorf("failed unmarshalling struct: %v", err)
	}
	header, err := q.header()
	if err != nil {
		return err
	}
	repoURL := q.base(repo)
	repoURL.Path = path.Join(repoURL.Path, subresource)
	req, _ := http.NewRequestWithContext(ctx, method, repoURL.String(), bytes.NewBuffer(content))
	req.Header = header
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("error performing rest call: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to download %s: %v: %v ", req.URL.String(), fmt.Errorf("status : %v", resp.StatusCode), string(body))
	}
	return nil
}

func (q *quayClient) Update(ctx context.Context, repository, description string) error {
	if err := q.json(ctx, http.MethodPut, repository, "", &Description{Description: description}); err != nil {
		return fmt.Errorf("error updating the repository description: %v", err)
	}

	if err := q.json(ctx, http.MethodPost, repository, "changevisibility", &Visibility{Visibility: "public"}); err != nil {
		return fmt.Errorf("error updating the repository visibility: %v", err)
	}
	return nil
}

func NewQuayClient(tokenFile, org string) *quayClient {
	return &quayClient{tokenFile: tokenFile, org: org}
}

type Description struct {
	Description string `json:"description"`
}

type Visibility struct {
	Visibility string `json:"visibility"`
}
