package almalinux

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"sort"
	"strings"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

//nolint:lll
const description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/1/13/AlmaLinux_Icon_Logo.svg/64px-AlmaLinux_Icon_Logo.svg.png" alt="drawing" height="15"/> AlmaLinux OS Generic Cloud images for KubeVirt.
<br />
<br />
Visit [almalinux.org](https://almalinux.org/) to learn more about the AlmaLinux OS project.`

type almalinux struct {
	Version         string
	Variant         string
	getter          http.Getter
	Arch            string
	ExampleUserData *docs.UserData
	EnvVariables    map[string]string
}

func (a *almalinux) Metadata() *api.Metadata {
	metadata := &api.Metadata{
		Name:         "almalinux",
		Version:      a.Version,
		Description:  description,
		EnvVariables: a.EnvVariables,
		Arch:         a.Arch,
	}

	if a.ExampleUserData != nil {
		metadata.ExampleUserData = *a.ExampleUserData
	}

	return metadata
}

func (a *almalinux) Inspect() (*api.ArtifactDetails, error) {
	baseURL := fmt.Sprintf("https://repo.almalinux.org/almalinux/%s/cloud/%s/images/", a.Version, a.Arch)
	checksumURL := baseURL + "CHECKSUM"

	raw, err := a.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the AlmaLinux checksum file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), hashsum.ChecksumFormatGNU)
	if err != nil {
		return nil, fmt.Errorf("error reading the AlmaLinux checksum file: %v", err)
	}

	pattern := regexp.MustCompile(
		fmt.Sprintf(`^AlmaLinux-%s-%s-\d+\.\d+-.+\.%s\.qcow2$`, a.Version, a.Variant, a.Arch),
	)
	candidates := []string{}
	for fileName := range checksums {
		if pattern.MatchString(fileName) {
			candidates = append(candidates, fileName)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates for version %q and variant %q found", a.Version, a.Variant)
	}

	sort.Strings(candidates)
	candidate := candidates[len(candidates)-1]

	prefix := fmt.Sprintf("AlmaLinux-%s-%s-", a.Version, a.Variant)
	suffix := fmt.Sprintf(".%s.qcow2", a.Arch)
	additionalTag := strings.TrimSuffix(strings.TrimPrefix(candidate, prefix), suffix)

	if checksum, exists := checksums[candidate]; exists {
		return &api.ArtifactDetails{
			Checksum:             checksum,
			ChecksumHash:         sha256.New,
			DownloadURL:          baseURL + candidate,
			AdditionalUniqueTags: []string{additionalTag},
			ImageArchitecture:    architecture.GetImageArchitecture(a.Arch),
		}, nil
	}

	return nil, fmt.Errorf("file %q does not exist in the sha256sum file: %v", a.Variant, err)
}

func (a *almalinux) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (a *almalinux) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (a *almalinux) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.GuestOsInfo,
		tests.SSH,
	}
}

func New(release, arch string, exampleUserData *docs.UserData, envVariables map[string]string) *almalinux {
	return &almalinux{
		Version:         release,
		Arch:            arch,
		Variant:         "GenericCloud",
		getter:          &http.HTTPGetter{},
		ExampleUserData: exampleUserData,
		EnvVariables:    envVariables,
	}
}
