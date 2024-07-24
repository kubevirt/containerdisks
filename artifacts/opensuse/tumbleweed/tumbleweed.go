package tumbleweed

import (
	"bytes"
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

type tumbleweed struct {
	Arch         string
	variant      string
	subVariant   string
	getter       http.Getter
	envVariables map[string]string
}

var _ api.Artifact = &tumbleweed{}

const description = `OpenSUSE Tumbleweed images for KubeVirt.
<br />
<br />
Visit [get.opensuse.org/tumbleweed/](https://get.opensuse.org/tumbleweed/) to learn more about OpenSUSE Tumbleweed.`

func (t *tumbleweed) Inspect() (*api.ArtifactDetails, error) {
	baseURL := "https://download.opensuse.org/tumbleweed/appliances/"
	raw, err := t.getter.GetAll(baseURL + "SHA256SUMS")
	if err != nil {
		return nil, fmt.Errorf("error downloading the tumbleweed SHA256SUMS file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), hashsum.ChecksumFormatGNU)
	if err != nil {
		return nil, fmt.Errorf("error reading the SHA256SUMS file: %v", err)
	}

	// openSUSE-Tumbleweed-Minimal-VM.x86_64-1.0.0-kvm-and-xen-Snapshot20240629.qcow2
	r := regexp.MustCompile(fmt.Sprintf(`%s\.%s-1\.0\.0-%s`, t.variant, t.Arch, t.subVariant))
	for file, checksum := range checksums {
		if r.MatchString(file) {
			return &api.ArtifactDetails{
				SHA256Sum:            checksum,
				DownloadURL:          baseURL + file,
				ImageArchitecture:    architecture.GetImageArchitecture(t.Arch),
				AdditionalUniqueTags: []string{"latest"},
			}, nil
		}
	}
	return nil, fmt.Errorf("variant %q does not exist in the SHA256SUMS file: %v", t.variant, err)
}

func (t *tumbleweed) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "opensuse-tumbleweed",
		Version:     "1.0.0",
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "opensuse",
		},
		EnvVariables: t.envVariables,
	}
}

func (t *tumbleweed) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (t *tumbleweed) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (t *tumbleweed) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.SSH,
	}
}

func New(arch string, envVariables map[string]string) *tumbleweed {
	return &tumbleweed{
		Arch:         arch,
		variant:      "openSUSE-Tumbleweed-Minimal-VM",
		subVariant:   "Cloud",
		getter:       &http.HTTPGetter{},
		envVariables: envVariables,
	}
}
