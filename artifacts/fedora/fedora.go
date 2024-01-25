package fedora

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
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
	Arch    string
	Variant string
	getter  http.Getter
}

type fedoraGatherer struct {
	Arch    string
	Variant string
	getter  http.Getter
}

const minimumVersion = 35

//nolint:lll
const description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Fedora_logo.svg/240px-Fedora_logo.svg.png" alt="drawing" width="15"/> Fedora [Cloud](https://alt.fedoraproject.org/cloud/) images for KubeVirt.
<br />
<br />
Visit [getfedora.org](https://getfedora.org/) to learn more about the Fedora project.`

func (f *fedora) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "fedora",
		Version:     f.Version,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "fedora",
		},
	}
}

func (f *fedora) Inspect() (*api.ArtifactDetails, error) {
	releases, err := getReleases(f.getter)
	if err != nil {
		return nil, fmt.Errorf("error getting releases: %v", err)
	}

	for i, release := range releases {
		if f.releaseMatches(&releases[i]) {
			components := strings.Split(release.Link, "/")
			fileName := components[len(components)-1]
			additionalTag := strings.TrimSuffix(strings.TrimPrefix(fileName, "Fedora-Cloud-Base-"), ".x86_64.qcow2")

			return &api.ArtifactDetails{
				SHA256Sum:            release.Sha256,
				DownloadURL:          release.Link,
				AdditionalUniqueTags: []string{additionalTag},
			}, nil
		}
	}

	return nil, fmt.Errorf("no release information in releases.json for fedora:%q found", f.Version)
}

func (f *fedora) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
		docs.WithSecureBoot(),
	)
}

func (f *fedora) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (f *fedora) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.GuestOsInfo,
		tests.SSH,
	}
}

func (f *fedoraGatherer) Gather() ([]api.Artifact, error) {
	releases, err := getReleases(f.getter)
	if err != nil {
		return nil, fmt.Errorf("error getting releases: %v", err)
	}

	artifacts := []api.Artifact{}
	for i, release := range releases {
		if f.releaseMatches(&releases[i]) {
			artifacts = append(artifacts, New(release.Version))
		}
	}

	return artifacts, nil
}

func getReleases(getter http.Getter) (Releases, error) {
	raw, err := getter.GetAll("https://getfedora.org/releases.json")
	if err != nil {
		return nil, fmt.Errorf("error downloading the fedora releases.json file: %v", err)
	}

	releases := Releases{}
	if err := json.Unmarshal(raw, &releases); err != nil {
		return nil, fmt.Errorf("error parsing the releases.json file: %v", err)
	}

	return releases, nil
}

func (f *fedora) releaseMatches(release *Release) bool {
	return release.Version == f.Version &&
		release.Arch == f.Arch &&
		release.Variant == f.Variant &&
		strings.HasSuffix(release.Link, "qcow2")
}

func (f *fedoraGatherer) releaseMatches(release *Release) bool {
	version, err := strconv.Atoi(release.Version)
	return err == nil && version >= minimumVersion &&
		release.Arch == f.Arch &&
		release.Variant == f.Variant &&
		strings.HasSuffix(release.Link, "qcow2")
}

func New(release string) *fedora {
	return &fedora{
		Version: release,
		Arch:    "x86_64",
		Variant: "Cloud",
		getter:  &http.HTTPGetter{},
	}
}

func NewGatherer() *fedoraGatherer {
	return &fedoraGatherer{
		Arch:    "x86_64",
		Variant: "Cloud",
		getter:  &http.HTTPGetter{},
	}
}
