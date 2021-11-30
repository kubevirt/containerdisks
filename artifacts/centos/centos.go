package centos

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
)

type centos struct {
	Version string
	Variant string
	getter  http.Getter
	Arch    string
}

func (c *centos) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:    "centos",
		Version: c.Version,
	}
}

func (c *centos) Inspect() (*api.ArtifactDetails, error) {
	var baseURL string
	var checksumURL string
	var checksumFormat hashsum.ChecksumFormat
	if strings.HasPrefix(c.Version, "8.") {
		baseURL = "https://cloud.centos.org/centos/8/x86_64/images/"
		checksumURL = baseURL + "CHECKSUM"
		checksumFormat = hashsum.ChecksumFormatBSD
	} else if strings.HasPrefix(c.Version, "7-") {
		baseURL = "https://cloud.centos.org/centos/7/images/"
		checksumURL = baseURL + "sha256sum.txt"
		checksumFormat = hashsum.ChecksumFormatGNU
	} else {
		panic(fmt.Sprintf("can't understand provided version: %q", c.Version))
	}
	raw, err := c.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the centos checksum file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), checksumFormat)
	if err != nil {
		return nil, fmt.Errorf("error reading the centos checksum file: %v", err)
	}

	candidates := []string{}
	if strings.HasPrefix(c.Version, "8.") {
		for fileName := range checksums {
			if strings.HasPrefix(fileName, fmt.Sprintf("CentOS-8-%s-%s", c.Variant, c.Version)) && strings.HasSuffix(fileName, "qcow2") {
				candidates = append(candidates, fileName)
			}
		}
	} else if strings.HasPrefix(c.Version, "7-") {
		components := strings.Split(c.Version, "-")
		for fileName := range checksums {
			if strings.HasPrefix(fileName, fmt.Sprintf("CentOS-7-x86_64-%s-%s.qcow2", c.Variant, components[1])) && strings.HasSuffix(fileName, "qcow2") {
				candidates = append(candidates, fileName)
			}
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates for version %q and variant %q found", c.Version, c.Variant)
	}

	sort.Strings(candidates)
	candidate := candidates[len(candidates)-1]

	var additionalTags []string
	if strings.HasPrefix(c.Version, "8.") {
		additionalTag := strings.TrimSuffix(strings.TrimPrefix(candidate, fmt.Sprintf("CentOS-8-%s-", c.Variant)), ".x86_64.qcow2")
		additionalTags = append(additionalTags, additionalTag)
		split := strings.Split(additionalTag, "-")
		if len(split) == 2 {
			additionalTags = append(additionalTags, split[0])
		}
	}

	if checksum, exists := checksums[candidate]; exists {
		return &api.ArtifactDetails{
			SHA256Sum:            checksum,
			DownloadURL:          baseURL + candidate,
			AdditionalUniqueTags: additionalTags,
		}, nil
	}
	return nil, fmt.Errorf("file %q does not exist in the sha256sum file: %v", c.Variant, err)

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
