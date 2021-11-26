package http

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
)

type Getter interface {
	GetAll(fileURL string) ([]byte, error)
	GetWithChecksum(fileURL string) (ReadCloserWithChecksum, error)
}

type ReadCloserWithChecksum interface {
	io.ReadCloser
	Checksum() string
}

type HTTPGetter struct {
}

func (H *HTTPGetter) GetAll(fileURL string) ([]byte, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to load promary repository file from %s: %v", fileURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Failed to download %s: %v ", fileURL, fmt.Errorf("status : %v", resp.StatusCode))
	}
	return ioutil.ReadAll(resp.Body)
}

func (H *HTTPGetter) GetWithChecksum(fileURL string) (ReadCloserWithChecksum, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to load promary repository file from %s: %v", fileURL, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("Failed to download %s: %v ", fileURL, fmt.Errorf("status : %v", resp.StatusCode))
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
