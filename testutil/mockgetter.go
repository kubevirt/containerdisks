package testutil

import (
	"io/ioutil"

	"kubevirt.io/containerdisks/pkg/http"
)

type mockGetter struct {
	mockFile string
}

func (m *mockGetter) GetAll(_ string) ([]byte, error) {
	return ioutil.ReadFile(m.mockFile)
}

func (m *mockGetter) GetWithChecksum(_ string) (http.ReadCloserWithChecksum, error) {
	panic("implement me")
}

func NewMockGetter(mockFile string) *mockGetter {
	return &mockGetter{mockFile: mockFile}
}
