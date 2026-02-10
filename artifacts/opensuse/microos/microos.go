package microos

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

type microos struct {
	Arch         string
	variant      string
	getter       http.Getter
	envVariables map[string]string
}

var _ api.Artifact = &microos{}

const description = `openSUSE MicroOS images for KubeVirt.
<br />
<br />
Visit [get.opensuse.org/microos/](https://get.opensuse.org/microos/) to learn more about openSUSE MicroOS.`

const (
	s390xArch           = "s390x"
	microOSVersion      = "16.0.0"
	microOSVersionRegex = `16\.0\.0`
)

func (t *microos) Inspect() (*api.ArtifactDetails, error) {
	baseURL := t.retrieveBaseURL()
	raw, err := t.getter.GetAll(baseURL + "SHA256SUMS")
	if err != nil {
		return nil, fmt.Errorf("error downloading the SHA256SUMS file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), hashsum.ChecksumFormatGNU)
	if err != nil {
		return nil, fmt.Errorf("error reading the SHA256SUMS file: %v", err)
	}

	// openSUSE-MicroOS.x86_64-16.0.0-OpenStack-Cloud-Snapshot20260207.qcow2
	// openSUSE-MicroOS.s390x-16.0.0-s390x-Cloud-Snapshot20260209.qcow2
	r := regexp.MustCompile(fmt.Sprintf(`%s\.%s-%s-%s`, t.variant, t.Arch, t.retrieveRegexpVersion(), t.subvariantByArchitecture()))
	for file, checksum := range checksums {
		if r.MatchString(file) {
			return &api.ArtifactDetails{
				Checksum:          checksum,
				ChecksumHash:      sha256.New,
				DownloadURL:       baseURL + file,
				ImageArchitecture: architecture.GetImageArchitecture(t.Arch),
			}, nil
		}
	}
	return nil, fmt.Errorf("variant %q does not exist in the SHA256SUMS file: %v", t.variant, err)
}

func (t *microos) retrieveBaseURL() string {
	if t.Arch == s390xArch {
		return "https://download.opensuse.org/ports/zsystems/tumbleweed/appliances/"
	}
	return "https://download.opensuse.org/tumbleweed/appliances/"
}

func (t *microos) subvariantByArchitecture() string {
	if t.Arch == s390xArch {
		return "s390x-Cloud"
	}
	return "OpenStack-Cloud"
}

func (t *microos) retrieveRegexpVersion() string {
	return microOSVersionRegex
}

func (t *microos) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "opensuse-microos",
		Version:     microOSVersion,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "opensuse",
		},
		EnvVariables: t.envVariables,
		Arch:         t.Arch,
	}
}

func (t *microos) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (t *microos) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (t *microos) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.SSH,
	}
}

func New(arch string, envVariables map[string]string) *microos {
	return &microos{
		Arch:         arch,
		variant:      "openSUSE-MicroOS",
		getter:       &http.HTTPGetter{},
		envVariables: envVariables,
	}
}
