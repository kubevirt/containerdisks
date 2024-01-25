package api

import (
	"context"
	"fmt"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/docs"
)

type ArtifactTest func(ctx context.Context, vmi *v1.VirtualMachineInstance, params *ArtifactTestParams) error

type ArtifactTestParams struct {
	// Username is the username used to log in into the VM.
	Username string
	// PrivateKey is the private key used to log in into the VM.
	PrivateKey interface{}
}

type ArtifactResult struct {
	// Tags contains all tags the built containerdisk was tagged with.
	Tags []string `json:",omitempty"`
	// Stage is the current stage of the containerdisk
	Stage string
	// Err indicates if an error happened while creating, verifying or promoting a containerdisk.
	Err string `json:",omitempty"`
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
	// Description of the project in Markdown format.
	Description string
	// CloudInit/Ignition Payload example.
	ExampleUserData docs.UserData
	// AdditionalLabels contains additional labels which should be added to the resulting containerdisk.
	// These labels can e.g. describe an appropriate instancetype or preference.
	AdditionalLabels map[string]string
}

func (m Metadata) Describe() string {
	return fmt.Sprintf("%s:%s", m.Name, m.Version)
}

type Artifact interface {
	Inspect() (*ArtifactDetails, error)
	Metadata() *Metadata
	VM(name, imgRef, userData string) *v1.VirtualMachine
	UserData(data *docs.UserData) string
	Tests() []ArtifactTest
}

type ArtifactsGatherer interface {
	// Gather must return a sorted list of dynamically gathered artifacts.
	// Artifacts have to be sorted in descending order with the latest release coming first.
	Gather() ([]Artifact, error)
}
