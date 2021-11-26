package api

import "fmt"

type ArtifactDetails struct {
	SHA256Sum   string
	DownloadURL string
}

type Metadata struct {
	Name    string
	Version string
}

func (m Metadata) Describe() string {
	return fmt.Sprintf("%s:%s", m.Name, m.Version)
}

type Artifact interface {
	Inspect() (*ArtifactDetails, error)
	Metadata() *Metadata
}
