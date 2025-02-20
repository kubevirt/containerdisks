package centosstream

import (
	"bytes"
	"fmt"
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
const description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/CentOS_Graphical_Symbol.svg/64px-CentOS_Graphical_Symbol.svg.png" alt="drawing" height="15"/> Centos Stream Generic Cloud images for KubeVirt.
<br />
<br />
Visit [centos.org](https://www.centos.org/) to learn more about the CentOS project.
<br />
Note that CentOS Stream 8 is EOL as of [May 31, 2024](https://blog.centos.org/2023/04/end-dates-are-coming-for-centos-stream-8-and-centos-linux-7/) and the associated containerdisks are now deprecated ahead of [removal in the future](https://github.com/kubevirt/containerdisks/issues/152).`

type centos struct {
	Version         string
	Variant         string
	getter          http.Getter
	Arch            string
	ExampleUserData *docs.UserData
	EnvVariables    map[string]string
}

func (c *centos) Metadata() *api.Metadata {
	metadata := &api.Metadata{
		Name:         "centos-stream",
		Version:      c.Version,
		Description:  description,
		EnvVariables: c.EnvVariables,
		Arch:         c.Arch,
	}

	if c.ExampleUserData != nil {
		metadata.ExampleUserData = *c.ExampleUserData
	}

	return metadata
}

func (c *centos) Inspect() (*api.ArtifactDetails, error) {
	var baseURL string

	if strings.HasPrefix(c.Version, "9") || strings.HasPrefix(c.Version, "10") {
		baseURL = fmt.Sprintf("https://cloud.centos.org/centos/%s-stream/%s/images/", c.Version, c.Arch)
	} else {
		panic(fmt.Sprintf("can't understand provided version: %q", c.Version))
	}

	checksumURL := baseURL + "CHECKSUM"
	checksumFormat := hashsum.ChecksumFormatBSD

	raw, err := c.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the centos stream checksum file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), checksumFormat)
	if err != nil {
		return nil, fmt.Errorf("error reading the centos stream checksum file: %v", err)
	}

	candidates := []string{}
	for fileName := range checksums {
		if strings.HasPrefix(fileName, fmt.Sprintf("CentOS-Stream-%s-%s", c.Variant, c.Version)) && strings.HasSuffix(fileName, "qcow2") {
			candidates = append(candidates, fileName)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates for version %q and variant %q found", c.Version, c.Variant)
	}

	sort.Strings(candidates)
	candidate := candidates[len(candidates)-1]

	var additionalTags []string
	suffix := fmt.Sprintf(".%s.qcow2", c.Arch)
	additionalTag := strings.TrimSuffix(strings.TrimPrefix(candidate, fmt.Sprintf("CentOS-Stream-%s-", c.Variant)), suffix)
	additionalTags = append(additionalTags, additionalTag)

	if checksum, exists := checksums[candidate]; exists {
		return &api.ArtifactDetails{
			SHA256Sum:            checksum,
			DownloadURL:          baseURL + candidate,
			AdditionalUniqueTags: additionalTags,
			ImageArchitecture:    architecture.GetImageArchitecture(c.Arch),
		}, nil
	}

	return nil, fmt.Errorf("file %q does not exist in the sha256sum file: %v", c.Variant, err)
}

func (c *centos) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (c *centos) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (c *centos) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.GuestOsInfo,
		tests.SSH,
	}
}

// New accepts CentOS Stream 8 and 9 versions.
func New(release, arch string, exampleUserData *docs.UserData, envVariables map[string]string) *centos {
	return &centos{
		Version:         release,
		Arch:            arch,
		Variant:         "GenericCloud",
		getter:          &http.HTTPGetter{},
		ExampleUserData: exampleUserData,
		EnvVariables:    envVariables,
	}
}
