//nolint:lll
package leap

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("OpenSUSE Leap", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(arch, version, mockFile string, envVariables map[string]string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(arch, version, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(details))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("leap:15.6 x86_64", "x86_64", "15.6", "testdata/openSUSE-Leap-15.6-Minimal-VM.x86_64-Cloud.qcow2.sha256",
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "opensuse.leap",
			},
			&api.ArtifactDetails{
				SHA256Sum:            "0f7f09a9a083088b51aa365fe0e4310e6b156c2153d6aa03a77b81eee884e52a",
				DownloadURL:          "https://download.opensuse.org/distribution/leap/15.6/appliances/openSUSE-Leap-15.6-Minimal-VM.x86_64-Cloud.qcow2",
				ImageArchitecture:    "amd64",
				AdditionalUniqueTags: []string{"15-latest"},
			},
			&api.Metadata{
				Name:        "opensuse-leap",
				Version:     "15.6",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "opensuse",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "opensuse.leap",
				},
			},
		),
		Entry("leap:15.6 aarch64", "aarch64", "15.6", "testdata/openSUSE-Leap-15.6-Minimal-VM.aarch64-Cloud.qcow2.sha256",
			nil,
			&api.ArtifactDetails{
				SHA256Sum:            "d2ff40176f8823ab869bf4d728f827ffd6c7f180940b9ccca865be6dc20b06dd",
				DownloadURL:          "https://download.opensuse.org/distribution/leap/15.6/appliances/openSUSE-Leap-15.6-Minimal-VM.aarch64-Cloud.qcow2",
				ImageArchitecture:    "arm64",
				AdditionalUniqueTags: []string{"15-latest"},
			},
			&api.Metadata{
				Name:        "opensuse-leap",
				Version:     "15.6",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "opensuse",
				},
			},
		),
	)
})

func TestLeap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpenSUSE Leap Suite")
}
