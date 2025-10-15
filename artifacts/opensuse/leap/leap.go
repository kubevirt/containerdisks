package leap

import (
	"crypto/sha256"
	"fmt"
	"strings"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

type leap struct {
	Arch         string
	Version      string
	getter       http.Getter
	envVariables map[string]string
}

var _ api.Artifact = &leap{}

const (
	leap16BaseURLFmt = "https://download.opensuse.org/distribution/leap/%s/appliances/Leap-%s-Minimal-VM.%s-Cloud.qcow2"
	leap15BaseURLFmt = "https://download.opensuse.org/distribution/leap/%s/appliances/openSUSE-Leap-%s-Minimal-VM.%s-Cloud.qcow2"
	description      = `OpenSUSE Leap images for KubeVirt.
<br />
<br />
Visit [get.opensuse.org/leap/](https://get.opensuse.org/leap/) to learn more about OpenSUSE Leap.`
)

func (l *leap) Inspect() (*api.ArtifactDetails, error) {
	baseURL := fmt.Sprintf(leap16BaseURLFmt, l.Version, l.Version, l.Arch)
	if strings.HasPrefix(l.Version, "15") {
		baseURL = fmt.Sprintf(leap15BaseURLFmt, l.Version, l.Version, l.Arch)
	}
	checksumBytes, err := l.getter.GetAll(baseURL + ".sha256")
	if err != nil {
		return nil, err
	}
	return &api.ArtifactDetails{
		Checksum:          strings.Split(string(checksumBytes), " ")[0],
		ChecksumHash:      sha256.New,
		DownloadURL:       baseURL,
		ImageArchitecture: architecture.GetImageArchitecture(l.Arch),
	}, nil
}

func (l *leap) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "opensuse-leap",
		Version:     l.Version,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "opensuse",
		},
		EnvVariables: l.envVariables,
		Arch:         l.Arch,
	}
}

func (l *leap) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (l *leap) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (l *leap) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.SSH,
	}
}

func New(arch, version string, envVariables map[string]string) *leap {
	return &leap{
		Arch:         arch,
		Version:      version,
		getter:       &http.HTTPGetter{},
		envVariables: envVariables,
	}
}
