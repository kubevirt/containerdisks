//nolint:lll
package microos

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"kubevirt.io/containerdisks/pkg/api"
	"kubevirt.io/containerdisks/pkg/common"
	"kubevirt.io/containerdisks/pkg/docs"
	"kubevirt.io/containerdisks/testutil"
)

var _ = Describe("openSUSE MicroOS", func() {
	DescribeTable("Inspect should be able to parse checksum files",
		func(arch, mockFile string, envVariables map[string]string, details *api.ArtifactDetails, metadata *api.Metadata) {
			c := New(arch, envVariables)
			c.getter = testutil.NewMockGetter(mockFile)
			got, err := c.Inspect()
			Expect(err).NotTo(HaveOccurred())
			Expect(got.ChecksumHash).ToNot(BeNil())
			Expect(got.Checksum).To(Equal(details.Checksum))
			Expect(got.DownloadURL).To(Equal(details.DownloadURL))
			Expect(got.AdditionalUniqueTags).To(Equal(details.AdditionalUniqueTags))
			Expect(got.ImageArchitecture).To(Equal(details.ImageArchitecture))
			Expect(got.Compression).To(Equal(details.Compression))
			Expect(c.Metadata()).To(Equal(metadata))
		},
		Entry("microos:1 x86_64", "x86_64", "testdata/microos.SHA256SUM",
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "opensuse.tumbleweed",
			},
			&api.ArtifactDetails{
				Checksum:          "bbc3613dfd22dac14d499afc44d1b67a8d6c8c2f77db71c4eb87887081104b7b",
				DownloadURL:       "https://download.opensuse.org/tumbleweed/appliances/openSUSE-MicroOS.x86_64-16.0.0-OpenStack-Cloud-Snapshot20260207.qcow2",
				ImageArchitecture: "amd64",
			},
			&api.Metadata{
				Name:        "opensuse-microos",
				Version:     "16.0.0",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "opensuse",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "opensuse.tumbleweed",
				},
				Arch: "x86_64",
			},
		),
		Entry("microos:1 s390x", "s390x", "testdata/microos-s390x.SHA256SUM",
			map[string]string{
				common.DefaultInstancetypeEnv: "u1.medium",
				common.DefaultPreferenceEnv:   "opensuse.tumbleweed",
			},
			&api.ArtifactDetails{
				Checksum:          "59d312f3f366ac9730343a27479f777319b591f93c4027c8702b7d794f123288",
				DownloadURL:       "https://download.opensuse.org/ports/zsystems/tumbleweed/appliances/openSUSE-MicroOS.s390x-16.0.0-s390x-Cloud-Snapshot20260207.qcow2",
				ImageArchitecture: "s390x",
			},
			&api.Metadata{
				Name:        "opensuse-microos",
				Version:     "16.0.0",
				Description: description,
				ExampleUserData: docs.UserData{
					Username: "opensuse",
				},
				EnvVariables: map[string]string{
					common.DefaultInstancetypeEnv: "u1.medium",
					common.DefaultPreferenceEnv:   "opensuse.tumbleweed",
				},
				Arch: "s390x",
			},
		),
	)
})

func TestMicroOS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "openSUSE MicroOS Suite")
}
