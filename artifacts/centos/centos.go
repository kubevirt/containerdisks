package centos

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
var description = `<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/CentOS_Graphical_Symbol.svg/64px-CentOS_Graphical_Symbol.svg.png" alt="drawing" height="15"/> Centos Generic Cloud images for KubeVirt.
<br />
<br />
Visit [centos.org](https://www.centos.org/) to learn more about the CentOS project.`

type centos struct {
	Version string
	Variant string
	getter  http.Getter
	Arch    string
}

func (c *centos) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:                   "centos",
		Version:                c.Version,
		Description:            description,
		ExampleUserDataPayload: c.UserData(&docs.UserData{}),
	}
}

func (c *centos) Inspect() (*api.ArtifactDetails, error) {
	baseURL, checksumURL, checksumFormat := getURLsAndChecksumFormat(c.Version)

	raw, err := c.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the centos checksum file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), checksumFormat)
	if err != nil {
		return nil, fmt.Errorf("error reading the centos checksum file: %v", err)
	}

	candidates := getCandidates(c.Version, c.Variant, checksums)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates for version %q and variant %q found", c.Version, c.Variant)
	}

	candidate := candidates[len(candidates)-1]
	if checksum, exists := checksums[candidate]; exists {
		return &api.ArtifactDetails{
			SHA256Sum:            checksum,
			DownloadURL:          baseURL + candidate,
			AdditionalUniqueTags: getAdditionalTags(c.Version, c.Variant, candidate),
		}, nil
	}

	return nil, fmt.Errorf("file %q does not exist in the sha256sum file: %v", c.Variant, err)
}

func getURLsAndChecksumFormat(version string) (baseURL string, checksumURL string, checksumFormat hashsum.ChecksumFormat) {
	switch {
	case strings.HasPrefix(version, "8."):
		baseURL = "https://cloud.centos.org/centos/8/x86_64/images/"
		checksumURL = baseURL + "CHECKSUM"
		checksumFormat = hashsum.ChecksumFormatBSD
	case strings.HasPrefix(version, "7-"):
		baseURL = "https://cloud.centos.org/centos/7/images/"
		checksumURL = baseURL + "sha256sum.txt"
		checksumFormat = hashsum.ChecksumFormatGNU
	default:
		panic(fmt.Sprintf("can't understand provided version: %q", version))
	}

	return
}

func getCandidates(version, variant string, checksums map[string]string) (candidates []string) {
	switch {
	case strings.HasPrefix(version, "8."):
		for fileName := range checksums {
			if strings.HasPrefix(fileName, fmt.Sprintf("CentOS-8-%s-%s", variant, version)) && strings.HasSuffix(fileName, "qcow2") {
				candidates = append(candidates, fileName)
			}
		}
	case strings.HasPrefix(version, "7-"):
		components := strings.Split(version, "-")
		for fileName := range checksums {
			if strings.HasPrefix(fileName, fmt.Sprintf("CentOS-7-x86_64-%s-%s.qcow2", variant, components[1])) &&
				strings.HasSuffix(fileName, "qcow2") {
				candidates = append(candidates, fileName)
			}
		}
	}

	sort.Strings(candidates)

	return
}

func getAdditionalTags(version, variant, candidate string) (additionalTags []string) {
	// The CentOS 8 version is expected to contain one dash
	const expectedCentos8VersionPartsCount = 2

	if strings.HasPrefix(version, "8.") {
		additionalTag := strings.TrimSuffix(strings.TrimPrefix(candidate, fmt.Sprintf("CentOS-8-%s-", variant)), ".x86_64.qcow2")
		additionalTags = append(additionalTags, additionalTag)
		split := strings.Split(additionalTag, "-")
		if len(split) == expectedCentos8VersionPartsCount {
			additionalTags = append(additionalTags, split[0])
		}
	}

	return
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

// New accepts CentOS 7 and 8 versions. Example patterns are 7-2111, 7-2009, 8.3, 8.4, ...
func New(release string) *centos {
	return &centos{
		Version: release,
		Arch:    "x86_64",
		Variant: "GenericCloud",
		getter:  &http.HTTPGetter{},
	}
}
