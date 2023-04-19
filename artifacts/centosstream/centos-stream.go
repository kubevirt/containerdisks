package centosstream

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

//nolint:lll
var description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/CentOS_Graphical_Symbol.svg/64px-CentOS_Graphical_Symbol.svg.png" alt="drawing" height="15"/> Centos Stream Generic Cloud images for KubeVirt.
<br />
<br />
Visit [centos.org](https://www.centos.org/) to learn more about the CentOS project.`

type centos struct {
	Version         string
	Variant         string
	getter          http.Getter
	Arch            string
	ExampleUserData *docs.UserData
}

func (c *centos) Metadata() *api.Metadata {
	metadata := &api.Metadata{
		Name:        "centos-stream",
		Version:     c.Version,
		Description: description,
	}

	if c.ExampleUserData != nil {
		metadata.ExampleUserData = *c.ExampleUserData
	}

	return metadata
}

func (c *centos) Inspect() (*api.ArtifactDetails, error) {
	var baseURL string

	if strings.HasPrefix(c.Version, "8") || strings.HasPrefix(c.Version, "9") {
		baseURL = fmt.Sprintf("https://cloud.centos.org/centos/%s-stream/x86_64/images/", c.Version)
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
	additionalTag := strings.TrimSuffix(strings.TrimPrefix(candidate, fmt.Sprintf("CentOS-Stream-%s-", c.Variant)), ".x86_64.qcow2")
	additionalTags = append(additionalTags, additionalTag)

	if checksum, exists := checksums[candidate]; exists {
		return &api.ArtifactDetails{
			SHA256Sum:            checksum,
			DownloadURL:          baseURL + candidate,
			AdditionalUniqueTags: additionalTags,
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
func New(release string, exampleUserData *docs.UserData) *centos {
	return &centos{
		Version:         release,
		Arch:            "x86_64",
		Variant:         "GenericCloud",
		getter:          &http.HTTPGetter{},
		ExampleUserData: exampleUserData,
	}
}
