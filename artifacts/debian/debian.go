package debian

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/architecture"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/pkg/http"
	"kubevirt.io/containerdisks/pkg/tests"
)

type Annotations struct {
	Digest string `json:"cloud.debian.org/digest"`
}

type ImageLabel struct {
	ImageFormat string `json:"upload.cloud.debian.org/image-format"`
	Version     string `json:"cloud.debian.org/version"`
}

type Metadata struct {
	Annotations Annotations `json:"annotations"`
	Labels      ImageLabel  `json:"labels"`
}

type Item struct {
	Metadata Metadata `json:"metadata"`
}

type BuildData struct {
	Items []Item `json:"items"`
}

type debian struct {
	Arch            string
	Version         string
	VersionName     string
	getter          http.Getter
	ExampleUserData *docs.UserData
	envVariables    map[string]string
}

const (
	baseURLFmt  = "https://cloud.debian.org/images/cloud/%s/latest/"
	baseNameFmt = "debian-%s-genericcloud-%s"
	description = `Debian Generic Cloud images for KubeVirt.
<br />
<br />
Visit [debian.org](https://cloud.debian.org/images/cloud/) to learn more about Debian project.`
)

var validDebianVersionPrefixes = []string{"11", "12", "13"}

func (d *debian) Metadata() *api.Metadata {
	metadata := &api.Metadata{
		Name:         "debian",
		Version:      d.Version,
		Description:  description,
		Arch:         d.Arch,
		EnvVariables: d.envVariables,
	}

	if d.ExampleUserData != nil {
		metadata.ExampleUserData = *d.ExampleUserData
	}

	return metadata
}

func decodeChecksum(base64Checksum string) (string, error) {
	base64Checksum = strings.TrimPrefix(base64Checksum, "sha512:")

	rawBytes, err := base64.RawStdEncoding.DecodeString(base64Checksum)
	if err != nil {
		return "", fmt.Errorf("error decoding debian digest: %v", err)
	}

	return hex.EncodeToString(rawBytes), nil
}

func (d *debian) getBuildData(jsonURL string) (additionalTags []string, checksum string, err error) {
	raw, err := d.getter.GetAll(jsonURL)
	if err != nil {
		return nil, "", fmt.Errorf("error downloading debian json file: %v", err)
	}

	var buildData BuildData
	if json.Unmarshal(raw, &buildData) != nil {
		return nil, "", fmt.Errorf("error decoding debian json file: %v", err)
	}

	if len(buildData.Items) == 0 {
		return nil, "", fmt.Errorf("build debian data not found")
	}

	for _, item := range buildData.Items {
		if item.Metadata.Labels.ImageFormat == "qcow2" {
			additionalTags = append(additionalTags, d.Version+"-"+item.Metadata.Labels.Version)
			checksum, err = decodeChecksum(item.Metadata.Annotations.Digest)
			return additionalTags, checksum, err
		}
	}

	return nil, "", fmt.Errorf("error locating the image information")
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func (d *debian) Inspect() (*api.ArtifactDetails, error) {
	if !hasAnyPrefix(d.Version, validDebianVersionPrefixes) {
		return nil, fmt.Errorf("can't understand provided version %s", d.Version)
	}

	baseURL := fmt.Sprintf(baseURLFmt, d.VersionName) + fmt.Sprintf(baseNameFmt, d.Version, architecture.GetImageArchitecture(d.Arch))

	additionalTags, checksum, err := d.getBuildData(baseURL + ".json")
	if err != nil {
		return nil, err
	}

	return &api.ArtifactDetails{
		Checksum:             checksum,
		ChecksumHash:         sha512.New,
		DownloadURL:          baseURL + ".qcow2",
		AdditionalUniqueTags: additionalTags,
		ImageArchitecture:    architecture.GetImageArchitecture(d.Arch),
	}, nil
}

func (d *debian) VM(name, imgRef, userData string) *v1.VirtualMachine {
	return docs.NewVM(
		name,
		imgRef,
		docs.WithRng(),
		docs.WithCloudInitNoCloud(userData),
	)
}

func (d *debian) UserData(data *docs.UserData) string {
	return docs.CloudInit(data)
}

func (d *debian) Tests() []api.ArtifactTest {
	return []api.ArtifactTest{
		tests.SSH,
	}
}

func New(version, versionName, arch string, exampleUserData *docs.UserData, envVariables map[string]string) *debian {
	return &debian{
		Version:         version,
		Arch:            arch,
		VersionName:     versionName,
		getter:          &http.HTTPGetter{},
		ExampleUserData: exampleUserData,
		envVariables:    envVariables,
	}
}
