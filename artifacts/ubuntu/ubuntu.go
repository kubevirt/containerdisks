package ubuntu

import (
	"bytes"
	"fmt"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

type ubuntu struct {
	Version      string
	Variant      string
	getter       http.Getter
	Arch         string
	Compression  string
	EnvVariables map[string]string
}

const description = `Ubuntu images for KubeVirt.
<br />
<br />
Visit [ubuntu.com](https://ubuntu.com/) to learn more about Ubuntu.`

func (u *ubuntu) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "ubuntu",
		Version:     u.Version,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "ubuntu",
		},
		EnvVariables: u.EnvVariables,
	}
}

func (u *ubuntu) Inspect() (*api.ArtifactDetails, error) {
	baseURL := fmt.Sprintf("https://cloud-images.ubuntu.com/releases/%v/release/", u.Version)
	checksumURL := baseURL + "SHA256SUMS"
	raw, err := u.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the ubuntu SHA256SUMS file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), hashsum.ChecksumFormatGNU)
	if err != nil {
		return nil, fmt.Errorf("error reading the SHA256SUMS.txt file: %v", err)
	}
	if checksum, exists := checksums[u.Variant]; exists {
		return &api.ArtifactDetails{
			SHA256Sum:         checksum,
			DownloadURL:       baseURL + u.Variant,
			Compression:       u.Compression,
			ImageArchitecture: architecture.GetImageArchitecture(u.Arch),
		}, nil
	}
	return nil, fmt.Errorf("file %q does not exist in the SHA256SUMS file: %v", u.Variant, err)
}

func (u *ubuntu) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (u *ubuntu) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (u *ubuntu) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.SSH,
	}
}

func New(release, arch string, envVariables map[string]string) *ubuntu {
	return &ubuntu{
		Version:      release,
		Arch:         arch,
		Variant:      fmt.Sprintf("ubuntu-%v-server-cloudimg-%s.img", release, architecture.GetImageArchitecture(arch)),
		getter:       &http.HTTPGetter{},
		EnvVariables: envVariables,
	}
}
