package api

import (
	"fmt"
)

type ArtifactDetails struct {
	// SHA256Sum is the checksum of the image to download.
	SHA256Sum string
	// DownloadURL points to the target image.
	DownloadURL string
	// Compression describes the compression format of the downloaded image.
	// Supported are "" (none), "gzip" and "xz".
	Compression string
	// AdditionalUniqueTags describes additional tags which furter specify the downloaded
	// artifact version. For instance the main moving tag for fedora 35 would be '35' and here additional tags
	// like '35-1.2'. This is useful for people to easier cross-reference the sources.
	AdditionalUniqueTags []string
}

type Metadata struct {
	// Name of the resulting container image in the remote container registry. For example "fedora".
	Name string
	// Version is the moving tag on the container image. For example "35".
	Version string

	// Description of the project in Markdown format
	Description string

	// CloudInit Payload example
	ExampleCloudInitPayload string
}

func (m Metadata) Describe() string {
	return fmt.Sprintf("%s:%s", m.Name, m.Version)
}

type Artifact interface {
	Inspect() (*ArtifactDetails, error)
	Metadata() *Metadata
}
