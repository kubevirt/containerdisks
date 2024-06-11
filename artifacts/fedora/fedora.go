package fedora

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/common"
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
	Version      string
	Arch         string
	Variant      string
	getter       http.Getter
	EnvVariables map[string]string
}

type fedoraGatherer struct {
	Version    string
	Archs      []string
	Variant    string
	Subvariant string
	getter     http.Getter
}

const minimumVersion = 38

//nolint:lll
const description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Fedora_logo.svg/240px-Fedora_logo.svg.png" alt="drawing" width="15"/> Fedora [Cloud](https://alt.fedoraproject.org/cloud/) images for KubeVirt.
<br />
<br />
Visit [getfedora.org](https://getfedora.org/) to learn more about the Fedora project.`

var additionalUniqueTagRegExp = regexp.MustCompile(`\d+-\d+\.\d+`)

func (f *fedora) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "fedora",
		Version:     f.Version,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "fedora",
		},
		EnvVariables: f.EnvVariables,
	}
}

func (f *fedora) Inspect() (*api.ArtifactDetails, error) {
	releases, err := getReleases(f.getter)
	if err != nil {
		return nil, fmt.Errorf("error getting releases: %v", err)
	}

	for i, release := range releases {
		if !f.releaseMatches(&releases[i]) {
			continue
		}

		details := &api.ArtifactDetails{
			SHA256Sum:         release.Sha256,
			DownloadURL:       release.Link,
			ImageArchitecture: architecture.GetImageArchitecture(f.Arch),
		}

		components := strings.Split(release.Link, "/")
		fileName := components[len(components)-1]
		if matches := additionalUniqueTagRegExp.FindStringSubmatch(fileName); len(matches) > 0 {
			details.AdditionalUniqueTags = append(details.AdditionalUniqueTags, matches[0])
		}

		return details, nil
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

func (f *fedoraGatherer) Gather() ([][]api.Artifact, error) {
	releases, err := getReleases(f.getter)
	if err != nil {
		return nil, fmt.Errorf("error getting releases: %v", err)
	}

	versions := map[string][]Release{}
	for i := range releases {
		release := releases[i]
		if f.releaseMatches(&release) {
			versions[release.Version] = append(versions[release.Version], release)
		}
	}

	// Ensure versions are always sorted with the latest first
	versionKeys := make([]string, 0, len(versions))
	for key := range versions {
		versionKeys = append(versionKeys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versionKeys)))

	var artifacts [][]api.Artifact
	for _, key := range versionKeys {
		releases := versions[key]
		var releaseArtifcats []api.Artifact
		for _, release := range releases {
			artifact := New(release.Version, release.Arch)
			releaseArtifcats = append(releaseArtifcats, artifact)
		}
		artifacts = append(artifacts, releaseArtifcats)
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
	if err != nil {
		return false
	}

	for _, arch := range f.Archs {
		if release.Arch == arch {
			return version >= minimumVersion &&
				release.Variant == f.Variant &&
				release.Subvariant == f.Subvariant &&
				strings.HasSuffix(release.Link, "qcow2")
		}
	}

	return false
}

const (
	defaultInstancetypeX86_64 = "u1.medium"
	defaultPreferenceX86_64   = "fedora"
)

func (f *fedora) setEnvVariables() {
	if f.Arch == "x86_64" {
		f.EnvVariables = map[string]string{
			common.DefaultInstancetypeEnv: defaultInstancetypeX86_64,
			common.DefaultPreferenceEnv:   defaultPreferenceX86_64,
		}
	}
}

func New(release, arch string) *fedora {
	f := &fedora{
		Version: release,
		Arch:    arch,
		Variant: "Cloud",
		getter:  &http.HTTPGetter{},
	}
	f.setEnvVariables()
	return f
}

func NewGatherer() *fedoraGatherer {
	return &fedoraGatherer{
		Archs:      []string{"x86_64", "aarch64"},
		Variant:    "Cloud",
		Subvariant: "Cloud_Base",
		getter:     &http.HTTPGetter{},
	}
}
