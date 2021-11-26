package fedora

import (
	"encoding/json"
	"fmt"
	"strings"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/http"
)

type Releases []Release

type Release struct {
	Subvariant string `json:"subvariant"`
	Variant    string `json:"variant"`
	Version    string `json:"version"`
	Link       string `json:"link"`
	Sha256     string `json:"sha256,omitempty"`
	Arch       string `json:"arch"`
	Size       string `json:"size,omitempty"`
}

type fedora struct {
	Version string
	getter  http.Getter
	Arch    string
	Variant string
}

func (f *fedora) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:    "fedora",
		Version: f.Version,
	}
}

func (f *fedora) Inspect() (*api.ArtifactDetails, error) {
	raw, err := f.getter.GetAll("https://getfedora.org/releases.json")
	if err != nil {
		return nil, fmt.Errorf("error downloading the fedora releases.json file: %v", err)
	}
	releases := Releases{}
	if err := json.Unmarshal(raw, &releases); err != nil {
		return nil, fmt.Errorf("error parsing the releases.json file: %v", err)
	}
	for _, release := range releases {
		if f.releaseMatches(&release) {
			return &api.ArtifactDetails{
				SHA256Sum:   release.Sha256,
				DownloadURL: release.Link,
			}, nil
		}
	}
	return nil, fmt.Errorf("no release information in releases.json for fedora:%q found", f.Version)
}

func (f *fedora) releaseMatches(release *Release) bool {
	return release.Version == f.Version &&
		release.Arch == f.Arch &&
		release.Variant == f.Variant &&
		strings.HasSuffix(release.Link, "qcow2")
}

func NewFedora(release string) api.Artifact {
	return &fedora{
		Version: release,
		Arch:    "x86_64",
		Variant: "Cloud",
		getter:  &http.HTTPGetter{},
	}
}
