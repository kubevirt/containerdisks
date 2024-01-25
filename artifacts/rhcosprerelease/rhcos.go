package rhcosprerelease

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/containers/image/v5/pkg/compression/types"
	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/hashsum"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

type rhcos struct {
	Version          string
	Variant          string
	getter           http.Getter
	Arch             string
	Compression      string
	AdditionalLabels map[string]string
}

//nolint:lll
const description = `RHCOS prerelease images for KubeVirt.
<br />
<br />
Visit [https://docs.openshift.com/container-platform/latest/architecture/architecture-rhcos.html) to learn more about Red Hat Enterprise Linux CoreOS.`

func (r *rhcos) Metadata() *api.Metadata {
	return &api.Metadata{
		Name:        "rhcos",
		Version:     strings.TrimPrefix(r.Version, "latest-") + "-pre-release",
		Description: description,
		ExampleUserData: docs.UserData{
			Username: "core",
		},
		AdditionalLabels: r.AdditionalLabels,
	}
}

func (r *rhcos) Inspect() (*api.ArtifactDetails, error) {
	baseURL := fmt.Sprintf("https://mirror.openshift.com/pub/openshift-v4/x86_64/dependencies/rhcos/pre-release/%s/", r.Version)
	checksumURL := baseURL + "sha256sum.txt"
	raw, err := r.getter.GetAll(checksumURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading the rhcos sha256sum.txt file: %v", err)
	}
	checksums, err := hashsum.Parse(bytes.NewReader(raw), hashsum.ChecksumFormatGNU)
	if err != nil {
		return nil, fmt.Errorf("error reading the sha256sum.txt file: %v", err)
	}

	var artifact *api.ArtifactDetails
	if checksum, exists := checksums[r.Variant]; exists {
		artifact = &api.ArtifactDetails{
			SHA256Sum:   checksum,
			DownloadURL: baseURL + r.Variant,
			Compression: r.Compression,
		}
	} else {
		return nil, fmt.Errorf("file %q does not exist in the sha256sum file: %v", r.Variant, err)
	}

	for variant, checksum := range checksums {
		if variant == r.Variant {
			continue
		}

		if checksum == artifact.SHA256Sum {
			additionalTag := strings.TrimSuffix(strings.TrimPrefix(variant, "rhcos-"), "-x86_64-openstack.x86_64.qcow2.gz")
			if !strings.Contains(additionalTag, "rc.") {
				continue
			}
			artifact.AdditionalUniqueTags = append(artifact.AdditionalUniqueTags, additionalTag)
		}
	}
	artifact.AdditionalUniqueTags = append(artifact.AdditionalUniqueTags, artifact.SHA256Sum)
	return artifact, nil
}

func (r *rhcos) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitConfigDrive(userData),
	)
}

func (r *rhcos) UserData(data *docs.UserData) string {
	return docs.Ignition(data)
}

func (r *rhcos) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.GuestOsInfo,
		tests.SSH,
	}
}

func New(release string, additionalLabels map[string]string) *rhcos {
	return &rhcos{
		Version:          release,
		Arch:             "x86_64",
		Variant:          "rhcos-openstack.x86_64.qcow2.gz",
		getter:           &http.HTTPGetter{},
		Compression:      types.GzipAlgorithmName,
		AdditionalLabels: additionalLabels,
	}
}
