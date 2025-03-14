package http

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
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
	Checksum() (string, string)
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
	sha256checksum := sha256.New()
	sha512checksum := sha512.New()

	// MultiWriter writes to both hashers simultaneously as the data is read through the Read method.
	// It allows to compute both shas in a single pass while the data will be still available to later read it and
	// create the containerdisk.
	teeReader := io.TeeReader(body, io.MultiWriter(sha256checksum, sha512checksum))

	return &readCloserWithChecksum{
		body:      body,
		teeReader: teeReader,
		sha256:    sha256checksum,
		sha512:    sha512checksum,
	}
}

type readCloserWithChecksum struct {
	body      io.ReadCloser
	teeReader io.Reader
	sha256    hash.Hash
	sha512    hash.Hash
}

func (r *readCloserWithChecksum) Read(p []byte) (n int, err error) {
	return r.teeReader.Read(p)
}

func (r *readCloserWithChecksum) Close() error {
	return r.body.Close()
}

func (r *readCloserWithChecksum) Checksum() (sha256checksum, sha512checksum string) {
	return hex.EncodeToString(r.sha256.Sum(nil)), hex.EncodeToString(r.sha512.Sum(nil))
}
