package http

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
)

type Getter interface {
	GetAll(fileURL string) ([]byte, error)
	GetAllWithContext(ctx context.Context, fileURL string) ([]byte, error)
	GetWithChecksum(fileURL string) (ReadCloserWithChecksum, error)
	GetWithChecksumAndContext(ctx context.Context, fileURL string) (ReadCloserWithChecksum, error)
}

type ReadCloserWithChecksum interface {
	io.ReadCloser
	Checksum() string
}

type HTTPGetter struct{}

func (h *HTTPGetter) GetAll(fileURL string) ([]byte, error) {
	return h.GetAllWithContext(context.Background(), fileURL)
}

func (h *HTTPGetter) GetAllWithContext(ctx context.Context, fileURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to load primary repository file from %s: %v", fileURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to load primary repository file from %s: %v", fileURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to download %s: %v ", fileURL, fmt.Errorf("status : %v", resp.StatusCode))
	}
	return io.ReadAll(resp.Body)
}

func (h *HTTPGetter) GetWithChecksum(fileURL string) (ReadCloserWithChecksum, error) {
	return h.GetWithChecksumAndContext(context.Background(), fileURL)
}

func (h *HTTPGetter) GetWithChecksumAndContext(ctx context.Context, fileURL string) (ReadCloserWithChecksum, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to load primary repository file from %s: %v", fileURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to load primary repository file from %s: %v", fileURL, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download %s: %v ", fileURL, fmt.Errorf("status : %v", resp.StatusCode))
	}
	return newReadCloserWithChecksum(resp.Body), nil
}

func newReadCloserWithChecksum(body io.ReadCloser) *readCloserWithChecksum {
	sha := sha256.New()
	teeReader := io.TeeReader(body, sha)
	return &readCloserWithChecksum{body: body, teeReader: teeReader, sha: sha}
}

type readCloserWithChecksum struct {
	body      io.ReadCloser
	teeReader io.Reader
	sha       hash.Hash
}

func (r *readCloserWithChecksum) Read(p []byte) (n int, err error) {
	return r.teeReader.Read(p)
}

func (r *readCloserWithChecksum) Close() error {
	return r.body.Close()
}

func (r *readCloserWithChecksum) Checksum() string {
	return hex.EncodeToString(r.sha.Sum(nil))
}
