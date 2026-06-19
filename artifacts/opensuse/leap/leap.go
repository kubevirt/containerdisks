package leap

import (
	"crypto/sha256"
	"fmt"
	"strconv"
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
	Username     string
	getter       http.Getter
	envVariables map[string]string
}

var _ api.Artifact = &leap{}

const (
	versionParts   = 2
	newURLMinMajor = 16
)

const description = `openSUSE Leap images for KubeVirt.
<br />
<br />
Visit [get.opensuse.org/leap/](https://get.opensuse.org/leap/) to learn more about openSUSE Leap.`

func (l *leap) Inspect() (*api.ArtifactDetails, error) {
	baseURL := l.buildBaseURL()
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

func (l *leap) buildBaseURL() string {
	parts := strings.SplitN(l.Version, ".", versionParts)
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < newURLMinMajor {
		return fmt.Sprintf(
			"https://download.opensuse.org/distribution/leap/%s/appliances/openSUSE-Leap-%s-Minimal-VM.%s-Cloud.qcow2",
			l.Version, l.Version, l.Arch,
		)
	}
	archPart := l.Arch
	if l.Arch == "s390x" {
		archPart = "s390x-s390x"
	}
	return fmt.Sprintf(
		"https://download.opensuse.org/distribution/leap/%s/appliances/Leap-%s-Minimal-VM.%s-Cloud.qcow2",
		l.Version, l.Version, archPart,
	)
}

func (l *leap) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "opensuse-leap",
		Version:     l.Version,
		Description: description,
		ExampleUserData: docs.UserData{
			Username: l.Username,
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

func New(arch, version, username string, envVariables map[string]string) *leap {
	return &leap{
		Arch:         arch,
		Version:      version,
		Username:     username,
		getter:       &http.HTTPGetter{},
		envVariables: envVariables,
	}
}
