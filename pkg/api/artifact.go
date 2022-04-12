package api

import (
	"fmt"

	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/containerdisks/pkg/docs"
)

type ArtifactResult struct {
	// Tags contains all tags the built containerdisk was tagged with.
	Tags []string
	// Verified indicates if the containerdisk was verified to be bootable and that the guest is working.
	Verified bool
}

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
	// CloudInit/Ignition Payload example
	ExampleUserDataPayload string
}

func (m Metadata) Describe() string {
	return fmt.Sprintf("%s:%s", m.Name, m.Version)
}

type Artifact interface {
	Inspect() (*ArtifactDetails, error)
	Metadata() *Metadata
	VM(name, imgRef, userData string) *v1.VirtualMachine
	UserData(data *docs.UserData) string
}
