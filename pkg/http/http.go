package http

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
)

type Getter interface {
	GetAll(fileURL string) ([]byte, error)
	GetAllWithContext(ctx context.Context, fileURL string) ([]byte, error)
	GetWithChecksum(fileURL string, checksumHasher func() hash.Hash) (ReadCloserWithChecksum, error)
	GetWithChecksumAndContext(ctx context.Context, fileURL string, checksumHasher func() hash.Hash) (ReadCloserWithChecksum, error)
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

func (h *HTTPGetter) GetWithChecksum(fileURL string, checksumHasher func() hash.Hash) (ReadCloserWithChecksum, error) {
	return h.GetWithChecksumAndContext(context.Background(), fileURL, checksumHasher)
}

func (h *HTTPGetter) GetWithChecksumAndContext(ctx context.Context, fileURL string, checksumHasher func() hash.Hash) (
	readCloser ReadCloserWithChecksum,
	err error,
) {
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
	return newReadCloserWithChecksum(resp.Body, checksumHasher), nil
}

func newReadCloserWithChecksum(body io.ReadCloser, checksumHasher func() hash.Hash) *readCloserWithChecksum {
	checksum := checksumHasher()
	teeReader := io.TeeReader(body, checksum)

	return &readCloserWithChecksum{
		body:         body,
		teeReader:    teeReader,
		checksumHash: checksum,
	}
}

type readCloserWithChecksum struct {
	body         io.ReadCloser
	teeReader    io.Reader
	checksumHash hash.Hash
}

func (r *readCloserWithChecksum) Read(p []byte) (n int, err error) {
	return r.teeReader.Read(p)
}

func (r *readCloserWithChecksum) Close() error {
	return r.body.Close()
}

func (r *readCloserWithChecksum) Checksum() string {
	return hex.EncodeToString(r.checksumHash.Sum(nil))
}
